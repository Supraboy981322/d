const { app, BrowserWindow, ipcMain } = require("electron");
const fs = require("fs");
const path = require("path");

function createWindow() {
  let ze_window = new BrowserWindow({
    width: 800,
    height: 600,
    autoHideMenuBar: true,
    webPreferences: {
      preload: path.join(__dirname, "preload.js")
    }
  });
  ze_window.loadFile("idx.html");
  ze_window.removeMenu();
}

app.whenReady().then (() => {
  ipcMain.handle("get_config", async () => {
    console.log("triggered get_config()");
    let conf_dir = path.join(app.getPath("home"), ".config", "Supraboy981322", "d");
    let conf_file = path.join(conf_dir, "d_electron.json");
    let conf;
    try {
      if (fs.existsSync(conf_file)) {
        let dat = fs.readFileSync(conf_file, "utf-8");
        conf = JSON.parse(dat);
        console.log("parsed config"); 
      } else {
        console.log("config not found");
        let def_conf = {
          server: "https://[your server address]"
        };
        fs.writeFileSync(conf_file, JSON.stringify(def_conf, null, 2), "utf-8");
        conf = def_conf;
      }
      return conf;
    } catch (e) {
      console.error(`failed to read config: ${e}`);
    }
  });

  ipcMain.handle("panic", (msg) => {
    console.error(msg);
    process.exit(1);
  });

  ipcMain.handle("quit", process.exit);

  ipcMain.handle("set_config", async (event, conf) => {
    console.log("triggered set_config()");
    let conf_dir = path.join(app.getPath("home"), ".config", "Supraboy981322", "d");
    let conf_file = path.join(conf_dir, "d_electron.json");
    try {
      fs.writeFileSync(conf_file, JSON.stringify(conf, null, 2), "utf-8");
      console.log("wrote config");
    } catch (e) {
      console.error(`failed to read config: ${e}`);
    }
  });
  app.on("activate", () => {
    if (BrowserWindow.getAllWindows().length === 0)
      createWindow();
  });
  createWindow();
});

app.on("window-all-closed", () => {
  app.quit(); //THIS ISN'T THE DEFAULT? (why?)
});
