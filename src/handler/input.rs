use std::time::Duration;

use crossterm::event::{self, Event, KeyCode, KeyEvent, MouseEvent, KeyModifiers};

#[derive(Debug)]
pub enum InputEvent {
    Keyboard(Vec<String>),
    Mouse(MouseEvent),
    Char(char),
}

pub fn keyboard_input_event_from_crossterm_key(event: KeyEvent) -> InputEvent {
    fn modifier_name(modifier: KeyModifiers) -> String {
        let dbg = format!("{:?}", modifier);

        dbg.trim_start_matches("KeyModifiers(")
            .trim_end_matches(')')
            .to_lowercase()
    }

    let mut keyboard_vec: Vec<String> = Vec::new();

    // Iterate all modifier flags
    let all_modifiers = [
        KeyModifiers::SHIFT,
        KeyModifiers::CONTROL,
        KeyModifiers::ALT,
        KeyModifiers::SUPER,
        KeyModifiers::HYPER,
        KeyModifiers::META,
    ];

    for modifier in all_modifiers {
        if event.modifiers.contains(modifier) {
            keyboard_vec.push(modifier_name(modifier));
        }
    }

    let key_str = match event.code {
        KeyCode::BackTab => String::from("tab"),
        KeyCode::F(n) => format!("f{}", n),
        KeyCode::Char(c) => return InputEvent::Char(c),
        other => format!("{:?}", other).to_lowercase(),
    };

    keyboard_vec.push(key_str);
    keyboard_vec.sort();

    InputEvent::Keyboard(keyboard_vec)
}

pub fn from_crossterm_event(event: crossterm::event::Event) -> Option<InputEvent> {
    match event {
        crossterm::event::Event::Key(key_event) => {
            let ev = keyboard_input_event_from_crossterm_key(key_event);
            Some(ev)
        },
        crossterm::event::Event::Mouse(mouse_event) => Some(InputEvent::Mouse(mouse_event)),
        _ => None,
    }
}