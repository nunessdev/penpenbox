package tui

type platformModel struct {
	platforms []platformItem // platforms in the list
	cursor    int            // which platform the cursor is pointing at
}
