const { contextBridge, ipcRenderer } = require('electron');

contextBridge.exposeInMainWorld('userAPI', {
  getUser: async () => {
    try {
      return await ipcRenderer.invoke('get-user');
    } catch (error) {
      console.error('Failed to fetch user:', error);
      throw error;
    }
  },
});

contextBridge.exposeInMainWorld('electron', {
  ipcRenderer: {
    invoke: (channel, ...args) => ipcRenderer.invoke(channel, ...args),
  },
});
