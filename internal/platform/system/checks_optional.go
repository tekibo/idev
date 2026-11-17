package system

func optionalSpecs() []spec {
	return []spec{
		{ID: "system.cmake", Packages: pkgs("cmake", "cmake", "cmake"), Command: "cmake"},
		{ID: "system.ninja", Packages: pkgs("ninja", "ninja-build", "ninja-build"), Command: "ninja"},
		{ID: "system.patchelf", Packages: pkgs("patchelf", "patchelf", "patchelf"), Command: "patchelf"},
		{ID: "system.fuse3", Packages: pkgs("fuse3", "libfuse3-dev", "fuse3-devel"), Pkg: "fuse3"},
	}
}
