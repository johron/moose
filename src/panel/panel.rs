use crate::handler::input::InputEvent;
use ratatui::layout::Rect;
use ratatui::Frame;
use std::fmt::Debug;

pub trait Panel: Debug {
    fn init(&mut self);
    fn is_initialized(&self) -> bool;
    fn identity(&self) -> &str;
    fn title(&self) -> String;
    fn render(&self, frame: &mut Frame, area: Rect);
    fn input(&mut self, input: InputEvent);
}