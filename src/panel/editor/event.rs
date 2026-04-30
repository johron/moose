use ropey::Rope;
use crate::panel::editor::editor::{Cursor, Editor, EditorEdit, EditorEvent};

pub fn handle_event(editor: &mut Editor, event: EditorEvent) {
    match event {
        EditorEvent::Mode(mode) => editor.mode = mode,
        EditorEvent::Edit(edit) => {
            match edit {
                EditorEdit::Insert(c) => insert_char(editor, c),
                EditorEdit::Delete => backspace(editor),
                _ => {}
            }
        },
        _ => {}
    }
}

fn cursor_to_char_idx(rope: &Rope, cursor: &Cursor) -> usize {
    let line_start = rope.line_to_char(cursor.line);
    line_start + cursor.col
}

fn insert_char(editor: &mut Editor, char: char) {
    let mut cursors: Vec<Cursor> = Vec::new();
    for cursor in editor.cursors.clone().into_iter() {
        editor.rope.insert_char(cursor_to_char_idx(&editor.rope, &cursor), char);
        cursors.push(Cursor::from(cursor.line, cursor.col.saturating_add(1), cursor.col.saturating_add(1)));
    }

    editor.cursors = cursors;
}

fn backspace(editor: &mut Editor) {
    let mut cursors: Vec<Cursor> = Vec::new();
    for cursor in editor.cursors.clone().into_iter() {
        editor.rope.remove(cursor_to_char_idx(&editor.rope, &cursor).saturating_sub(1)..cursor_to_char_idx(&editor.rope, &cursor));
        cursors.push(Cursor::from(cursor.line, cursor.col.saturating_sub(1), cursor.col.saturating_sub(1)));
    }

    editor.cursors = cursors;
}