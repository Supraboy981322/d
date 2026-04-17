const rl = @import("raylib");

pub fn next_key() ?rl.KeyboardKey {
    const k = rl.getKeyPressed();
    return if (k == .null) null else k;
}
