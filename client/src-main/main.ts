import { app, BrowserWindow, ipcMain, protocol, net } from 'electron';
import path from 'path';
import url from 'url';
import { stat } from 'node:fs/promises';

// Handle creating/removing shortcuts on Windows when installing/uninstalling.
import electronSquirrelStartup from 'electron-squirrel-startup';
if(electronSquirrelStartup) app.quit();

// Only one instance of the electron main process should be running due to how chromium works.
// If another instance of the main process is already running `app.requestSingleInstanceLock()`
// will return false, `app.quit()` will be called, and the other instances will receive a
// `'second-instance'` event.
// https://www.electronjs.org/docs/latest/api/app#apprequestsingleinstancelockadditionaldata
if(!app.requestSingleInstanceLock()) {
  app.quit();
}

// This event will be called when a second instance of the app tries to run.
// https://www.electronjs.org/docs/latest/api/app#event-second-instance
app.on('second-instance', (event, args, workingDirectory, additionalData) => {
  createWindow();
});

const scheme = 'app';
const srcFolder = path.join(app.getAppPath(), `.vite/main_window/`);
const staticAssetsFolder = import.meta.env.DEV ? path.join(import.meta.dirname, '../../static/') : srcFolder;

protocol.registerSchemesAsPrivileged([{
  scheme: scheme,
  privileges: {
    standard: true,
    secure: true,
    allowServiceWorkers: true,
    supportFetchAPI: true,
    corsEnabled: false,
  },
}]);

app.on('ready', () => {
  protocol.handle(scheme, async (request) => {
    const requestPath = path.normalize(decodeURIComponent(new URL(request.url).pathname));

    async function isFile(filePath: string) {
      try {
	if((await stat(filePath)).isFile()) return filePath;
      }
      catch(e) {}
    }

    const responseFilePath = await isFile(path.join(srcFolder, requestPath))
      ?? await isFile(path.join(srcFolder, path.dirname(requestPath), `${path.basename(requestPath) || 'index'}.html`))
      ?? path.join(srcFolder, '200.html');

    return await net.fetch(url.pathToFileURL(responseFilePath).toString());
  });
});

function createWindow() {
  const mainWindow = new BrowserWindow({
    width: 1400,
    height: 900,
    minWidth: 700,
    minHeight: 400
  });

  if(import.meta.env.DEV) {
    mainWindow.loadURL(VITE_DEV_SERVER_URLS['main_window']);

    mainWindow.webContents.openDevTools();
  }
  else {
    mainWindow.loadURL('app://-/');
  }
}

app.on('ready', createWindow);

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

// In this file you can include the rest of your app's specific main process
// code. You can also put them in separate files and import them here.
