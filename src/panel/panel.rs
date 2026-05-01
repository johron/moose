use crate::handler::input::InputEvent;
use ratatui::layout::Rect;
use ratatui::Frame;
use std::fmt::Debug;

pub trait Panel: Debug {
    fn init(&mut self);
    fn is_initialized(&self) -> bool;
    fn is_normal_mode(&self) -> bool;
    fn identity(&self) -> &str;
    fn title(&self) -> String;
    fn render(&self, frame: &mut Frame, area: Rect);
    fn input(&mut self, input: InputEvent);
}

#[derive(Debug, Clone, Eq, Hash, PartialEq)]
pub struct Cursor {
    pub line: usize,
    pub col: usize,
    pub goal: usize,
}

impl Cursor {
    pub fn new() -> Self {
        Cursor {
            line: 0,
            col: 0,
            goal: 0,
        }
    }

    pub fn from(line: usize, col: usize, goal: usize) -> Self {
        Cursor {
            line,
            col,
            goal,
        }
    }
}