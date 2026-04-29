use crossterm::event::MouseEvent;

pub enum InputEvent {
    Keyboard(Vec<String>),
    Mouse(MouseEvent),
    Char(char),
}

