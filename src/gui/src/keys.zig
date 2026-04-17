const std = @import("std");
const rl = @import("raylib");

pub const Key = struct {
    tag:rl.KeyboardKey,
    last_active:f64 = 0.0,
    repeat_delay:f64,

    pub fn init(
        key:rl.KeyboardKey,
        repeat_delay:f64,
    ) Key {
        return .{
            .tag = key,
            .repeat_delay = repeat_delay,
        };
    }

    pub fn is_active(self:*Key, t:f64) bool {
        if (self.last_active + self.repeat_delay <= t and rl.isKeyDown(self.tag)) {
            self.last_active = t;
            return true;
        }
        return false;
    }
};

pub fn is_active(key:rl.KeyboardKey, t:f64) bool {
    for (&keys) |*k|
        if (k.tag == key)
            return @constCast(k).is_active(t);
    unreachable;
}

const num_keys:comptime_int = std.meta.tags(rl.KeyboardKey).len;
pub const keys:[num_keys]Key = b: {
    var buf:[num_keys]Key = undefined;
    for (std.meta.tags(rl.KeyboardKey), 0..) |key, i| {
        buf[i] = .init(key, 0.05);
    }
    break :b buf;
};

pub const KeyItr = struct {
    buf:[]Key,
    pos:usize = 0,

    pub fn init(alloc:std.mem.Allocator) !KeyItr {
        const time = rl.getTime();
        return .{
            .buf = b: {
                var buf:[]Key = try alloc.alloc(Key, 0);
                for (&keys) |*key| if (@constCast(key).is_active(time)) {
                    const old_buf = try alloc.dupe(Key, buf);
                    defer alloc.free(old_buf);
                    alloc.free(buf);
                    buf = try alloc.alloc(Key, old_buf.len + 1);
                    buf[old_buf.len] = key.*;
                    std.debug.print("active: {s}\n", .{@tagName(key.tag)});
                };
                break :b buf;
            },
        };
    }

    pub fn deinit(self:*KeyItr, alloc:std.mem.Allocator) void {
        alloc.free(self.buf);
    }

    pub fn next(self:*KeyItr) ?Key {
        if (self.buf.len >= self.pos) return null;
        defer self.pos += 1;
        return self.buf[self.pos];
    }
};
