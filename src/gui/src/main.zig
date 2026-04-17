const std = @import("std");
const rl = @import("raylib");

pub fn main() !void {
    const screenWidth = 800;
    const screenHeight = 480;

    rl.initWindow(screenWidth, screenHeight, "d_gui");
    defer rl.closeWindow();

    if (rl.getMonitorCount() > 1)
        rl.setWindowMonitor(1);

    rl.setTargetFPS(30);

    while (!rl.windowShouldClose()) {
        rl.beginDrawing();
        defer rl.endDrawing();
        rl.clearBackground(.black);
        rl.drawText("foo", 100, 100, 20, .white);
    }
}
