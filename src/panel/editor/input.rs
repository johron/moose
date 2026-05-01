use crate::handler::config::{char_config, vec_config};
use crate::handler::input::InputEvent;
use crate::panel::editor::editor::{Editor, EditorEdit, EditorEvent, EditorMode};

pub fn handle_input(editor: &Editor, input_event: InputEvent) -> EditorEvent {
    match editor.mode {
        EditorMode::Normal => handle_normal_mode_input(editor, input_event),
        EditorMode::Insert => handle_insert_mode_input(editor, input_event),
    }
}

fn handle_normal_mode_input(editor: &Editor, input_event: InputEvent) -> EditorEvent {
    let config = editor.config.clone();

    match input_event {
        InputEvent::Char(c) => {
            if Some(c) == char_config(config.enter_insert_mode) {
                EditorEvent::Mode(EditorMode::Insert)
            } else {
                EditorEvent::None
            }
        },
        _ => EditorEvent::None,
    }
}

fn handle_insert_mode_input(editor: &Editor, input_event: InputEvent) -> EditorEvent {
    let config = editor.config.clone();

    match input_event {
        InputEvent::Char(c) => EditorEvent::Edit(EditorEdit::Insert(c)),
        InputEvent::Keyboard(vec) => {
            if vec == vec_config(config.enter_normal_mode) {
                EditorEvent::Mode(EditorMode::Normal)
            } else if vec == vec_config(config.insert_mode_delete_char) {
                EditorEvent::Edit(EditorEdit::Delete)
            } else {
                EditorEvent::None
            }
        },
        _ => EditorEvent::None,
    }
}
