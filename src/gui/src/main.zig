const std = @import("std");
const rl = @import("raylib");
const hlp = @import("helpers.zig");

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const alloc = gpa.allocator();

    rl.initWindow(800, 480, "d_gui");
    defer rl.closeWindow();

    if (rl.getMonitorCount() > 1)
        rl.setWindowMonitor(1);

    rl.setTargetFPS(30);

    var frame_count:usize = 0;
    //var letter_num:usize = 0;
    var input_txt = try std.ArrayList(u8).initCapacity(alloc, 0);
    defer input_txt.deinit(alloc);

    while (!rl.windowShouldClose()) {
        frame_count +%= 1;

        const txt_box:rl.Rectangle = b: {
            const screen_width = @as(f32, @floatFromInt(rl.getScreenWidth()));
            const width:f32 = screen_width * 0.9;
            break :b .{
                .x = screen_width * 0.05,
                .y = @as(f32, @floatFromInt(rl.getScreenHeight())) - 70,
                .width = width,
                .height = 35,
            };
        };

        while (hlp.next_key()) |key| switch (key) {
            .backspace => _ = input_txt.pop(),
            else => {
                const key_int:c_int = @intFromEnum(key);
                if (key_int >= 32 and key_int <= 125) {
                    try input_txt.append(alloc, @intCast(key_int));
                }
            },
        };

        rl.beginDrawing();
        defer rl.endDrawing();
        rl.clearBackground(.black);
        rl.drawText("foo", 100, 100, 20, .white);
    }
}
