use std::fmt::Debug;
use ratatui::Frame;
use ratatui::layout::Rect;
use crate::handler::input::InputEvent;

pub trait Panel: Debug {
    fn init(&mut self);
    fn is_init(&self) -> bool;
    fn identity(&self) -> &str;
    fn title(&self) -> String;
    fn render(&self, frame: &mut Frame, area: Rect);
    fn input(&mut self, input: InputEvent);
}