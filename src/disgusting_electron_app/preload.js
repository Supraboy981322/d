//this feels so needless
const { contextBridge, ipcRenderer } = require("electron");

contextBridge.exposeInMainWorld("api", {
  panic: () => ipcRenderer.invoke("panic"),
  quit: () => ipcRenderer.invoke("quit"),
  get_config: () => ipcRenderer.invoke("get_config"),
  set_config: (conf) => ipcRenderer.invoke("set_config", conf),
});
