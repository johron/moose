use crate::handler::input::InputEvent;
use crate::panel::editor::editor::{Editor, EditorEdit, EditorEvent, EditorMode};

pub fn handle_input(editor: &Editor, input_event: InputEvent) -> EditorEvent {
    match editor.mode {
        EditorMode::Normal => handle_normal_mode_input(input_event),
        EditorMode::Insert => handle_insert_mode_input(input_event),
        EditorMode::Command => handle_command_mode_input(input_event),
    }
}

fn handle_normal_mode_input(input_event: InputEvent) -> EditorEvent {
    match input_event {
        InputEvent::Char(c) => {
            match c {
                'i' => EditorEvent::Mode(EditorMode::Insert),
                _ => EditorEvent::None,
            }
        }
        _ => EditorEvent::None,
    }
}

fn handle_insert_mode_input(input_event: InputEvent) -> EditorEvent {
    match input_event {
        InputEvent::Char(c) => EditorEvent::Edit(EditorEdit::Insert(c)),
        InputEvent::Keyboard(vec) => {
            if vec == ["esc"] {
                EditorEvent::Mode(EditorMode::Normal)

            } else {
                EditorEvent::None
            }
        },
        _ => EditorEvent::None,
    }
}

fn handle_command_mode_input(input_event: InputEvent) -> EditorEvent {
    todo!()
}