const std = @import("std");
const rl = @import("raylib");
const rg = @import("raygui");
const hlp = @import("helpers.zig");

pub fn main(init:std.process.Init) !void {
    const alloc = init.gpa;

    rl.initWindow(800, 480, "d_gui");
    defer rl.closeWindow();

    if (rl.getMonitorCount() > 1)
        rl.setWindowMonitor(1);

    rl.setTargetFPS(30);

    while (!rl.windowShouldClose()) {
        {
            var itr = try @import("keys.zig").KeyItr.init(alloc);
            defer itr.deinit(alloc);
            while (itr.next()) |key| switch (key.tag) {
                else => {},
            };
        }

        rl.beginDrawing();
        defer rl.endDrawing();

        rl.clearBackground(.black);
        rl.drawText("foo", 200, 200, 20, .ray_white);
    }
}
