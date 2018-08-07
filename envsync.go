package envsync

// EnvSyncer describes some contracts to synchronize env.
type EnvSyncer interface {
	// Sync synchronizes source and target.
	// Source is the default env or the sample env.
	// Target is modified env.
	// Both source and target are string and indicate the location of the files.
	//
	// Any values in source that aren't in target will be written to target.
	// Any values in source that are in target won't be written to target.
	Sync(source, target string) error
}
