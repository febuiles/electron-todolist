import { app, net } from 'electron';
import path from 'path';
import { existsSync, readFileSync, writeFileSync } from 'fs';
import { AppHost } from "./config"

const userDataPath = app.getPath('userData');
export const userFilePath = path.join(userDataPath, 'user.json');

export async function createUser() {
  try {
    const response = await fetch(`${AppHost}/users/`, {
      method: 'POST'
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const user = await response.json();
    return user;
  } catch (err) {
    throw err;
  }
}

export async function initializeUser() {
  if (existsSync(userFilePath)) {
    const userData = readFileSync(userFilePath, 'utf-8');
    const user = JSON.parse(userData);
    return user;
  } else {
    try {
      const user = await createUser();
      // after creating a user, we need to create an initial todolist
      const tlRes = await fetch(`${AppHost}/todolists/`, {
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
