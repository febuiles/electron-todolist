import { app, net } from 'electron';
import path from 'path';
import { existsSync, readFileSync, writeFileSync } from 'fs';

const userDataPath = app.getPath('userData');
const userFilePath = path.join(userDataPath, 'user.json');
const app_host = "http://localhost:8080"

export async function createUser() {
  const request = net.request({
    method: 'POST',
    url: `${app_host}/users`,
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
  if (existsSync(userFilePath)) {
    const userData = readFileSync(userFilePath, 'utf-8');
    const user = JSON.parse(userData);
    return user;
  } else {
    try {
      const user = await createUser();
      console.log(userFilePath)
      writeFileSync(userFilePath, JSON.stringify(user));
      return user;
    } catch (err) {
      console.error('Error creating user:', err);
      throw err;
    }
  }
}
