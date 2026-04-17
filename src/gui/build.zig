const std = @import("std");

pub fn build(b: *std.Build) void {
    //build settings
    const bin = b.addExecutable(.{
        .name = "d_gui",
        .root_module = b.createModule(.{
            .root_source_file = b.path("src/main.zig"),
            .target = b.graph.host,
        }),
    });

    const raylib_dep = b.dependency("raylib_zig", .{
        .target = b.graph.host,
        .linkage = .static,
    });

    const raylib = raylib_dep.module("raylib"); // main raylib module
    const raygui = raylib_dep.module("raygui"); // raygui module
    const raylib_artifact = raylib_dep.artifact("raylib"); // raylib C library

    b.installArtifact(bin);
    bin.root_module.linkLibrary(raylib_artifact);
    bin.root_module.addImport("raylib", raylib);
    bin.root_module.addImport("raygui", raygui);

    //for 'zig build run'
    const run_bin = b.addRunArtifact(bin);
    if (b.args) |args| {
        run_bin.addArgs(args);
    }
    const run_step = b.step("run", "run the program");
    run_step.dependOn(&run_bin.step);
}
