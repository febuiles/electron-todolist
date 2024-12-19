import { app, net } from 'electron';
import path from 'path';
import { existsSync, readFileSync, writeFileSync } from 'fs';

const userDataPath = app.getPath('userData');
const userFilePath = path.join(userDataPath, 'user.json');
const app_host = "http://localhost:8080"

export async function createUser() {
  const request = net.request({
    method: 'POST',
    url: `${app_host}/users/`,
  });

  return new Promise((resolve, reject) => {
    let responseData = '';

    request.on('response', (response) => {
      response.on('data', (chunk) => {
        responseData += chunk.toString();
      });

      response.on('end', () => {
        try {
          const user = JSON.parse(responseData);
          resolve(user);
        } catch (err) {
          reject(err);
        }
      });
    });

    request.on('error', (err) => {
      reject(err);
    });

    request.end();
  });
}

export async function initializeUser() {
  console.log(userFilePath)
  if (existsSync(userFilePath)) {
    const userData = readFileSync(userFilePath, 'utf-8');
    const user = JSON.parse(userData);
    return user;
  } else {
    try {
      const user = await createUser();
      // assign an initial todolist
      const tlRes = await fetch(`http://localhost:8080/todolists/`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: user.id }),
      });

      const tl = await tlRes.json();

      user.lastUsedTodolistId = tl.id

      writeFileSync(userFilePath, JSON.stringify(user));
      return user;
    } catch (err) {
      console.error('Error creating user:', err);
      throw err;
    }
  }
}

export async function getUser(id: number) {
  const request = net.request({
    method: 'GET',
    url: `${app_host}/users/${id}`,
  });

  return new Promise((resolve, reject) => {
    let responseData = '';
    request.on('response', (response) => {
      if (response.statusCode === 404) {
        reject(new Error('User not found'));
        return;
      }

      response.on('data', (chunk) => {
        responseData += chunk.toString();
      });

      response.on('end', () => {

        try {
          const user = JSON.parse(responseData);
          resolve(user);
        } catch (err) {
          reject(err);
        }
      });
    });

    request.on('error', (err) => {
      console.log("aqui");

      console.error('Error getting user:', err);
      reject(err);
    });
    request.end();
  });
}
