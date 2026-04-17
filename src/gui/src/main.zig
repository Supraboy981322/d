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

        {
            var itr = try @import("keys.zig").KeyItr.init(alloc);
            defer itr.deinit(alloc);
            while (itr.next()) |key| switch (key.tag) {
                .backspace => _ = input_txt.pop(),
                else => {
                    const key_int:c_int = @intFromEnum(key.tag);
                    if (key_int >= 32 and key_int <= 125) {
                        try input_txt.append(alloc, @intCast(key_int));
                    }
                },
            };
        }

        rl.beginDrawing();
        defer rl.endDrawing();

        rl.clearBackground(.black);
        rl.drawRectangleRec(txt_box, .white);

        const elderly_input_buf:[:0]const u8 = b: {
            var itr = std.mem.reverseIterator(input_txt.items);
            var buf = try std.ArrayList(u8).initCapacity(alloc, 0);
            defer buf.deinit(alloc);
            var buf_elderly = try alloc.allocSentinel(u8, 0, 0);
            while (itr.next()) |b| {
                try buf.append(alloc, b);
                alloc.free(buf_elderly);
                buf_elderly = try alloc.dupeZ(u8, buf.items);
                if (@as(f32, @floatFromInt(rl.measureText(buf_elderly, 20))) > txt_box.width) break;
            }
            for (0..3) |_| _ = buf.pop();
            alloc.free(buf_elderly);
            buf_elderly = try alloc.dupeZ(u8, buf.items);
            std.mem.reverse(u8, buf_elderly);
            break :b buf_elderly;
        };
        defer alloc.free(elderly_input_buf);

        rl.drawText(
            elderly_input_buf,
            @intFromFloat(txt_box.x + 5),
            @intFromFloat(txt_box.y + 8),
            20,
            .dark_gray,
        );

        if ((frame_count / 20) % 2 == 0) rl.drawText(
            "_",
            @as(i32, @intFromFloat(txt_box.x)) + 8 + rl.measureText(elderly_input_buf, 20),
            @intFromFloat(txt_box.y + 12),
            20,
            .blue
        );
    }
}
