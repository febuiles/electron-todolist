import { app, BrowserWindow, protocol, net, ipcMain } from 'electron';
import path from 'path';
import url from 'url';

import { createUser, initializeUser } from './users'

const scheme = 'app';
const srcFolder = path.join(app.getAppPath(), `.vite/main_window/`);

protocol.registerSchemesAsPrivileged([
  {
    scheme: scheme,
    privileges: {
      standard: true,
      secure: true,
      allowServiceWorkers: true,
      supportFetchAPI: true,
      corsEnabled: false,
    },
  },
]);

app.on('ready', async () => {
  try {

    const user = await initializeUser();
    ipcMain.handle('get-user', () => { // we'll need this in the renderer
      return user;
    });

    protocol.handle(scheme, async (request) => {
      const requestPath = path.normalize(decodeURIComponent(new URL(request.url).pathname));

      const responseFilePath = path.join(srcFolder, requestPath) || path.join(srcFolder, '200.html');
      return await net.fetch(url.pathToFileURL(responseFilePath).toString());
    });

    createWindow();
  } catch (err) {
    console.error('Failed to initialize user:', err);
    app.quit();
  }
});

function createWindow() {
  const mainWindow = new BrowserWindow({
    width: 1400,
    height: 900,
    minWidth: 700,
    minHeight: 400,
    webPreferences: {
      preload: path.join(import.meta.dirname, '../../renderer/preload.js'),
    },
  });

  if (import.meta.env.DEV) {
    mainWindow.loadURL(VITE_DEV_SERVER_URLS['main_window']);
    mainWindow.webContents.openDevTools();
  } else {
    mainWindow.loadURL('app://-/');
  }
}

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});
