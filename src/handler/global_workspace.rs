use ratatui::Frame;
use ratatui::layout::{Constraint, Direction, Layout, Rect};
use crate::handler::renderer::render;
use crate::handler::workspace::Workspace;
use crate::panel::global_panel::GlobalPanel;
use crate::panel::panel::Panel;

#[derive(Debug)]
pub struct GlobalWorkspace {
    pub(crate) bottom: Option<Box<dyn GlobalPanel>>,   // 0
    pub(crate) left: Option<Box<dyn GlobalPanel>>,     // 1
    pub(crate) right: Option<Box<dyn GlobalPanel>>,    // 2
    active: usize, //                          ^
}

impl GlobalWorkspace {
    pub fn new() -> Self {
        GlobalWorkspace {
            bottom: None,
            left: None,
            right: None,
            active: 0,
        }
    }

    pub fn set_bottom(&mut self, panel: Box<dyn GlobalPanel>) {
        self.bottom = Some(panel);
    }

    pub fn set_left(&mut self, panel: Box<dyn GlobalPanel>) {
        self.left = Some(panel);
    }

    pub fn set_right(&mut self, panel: Box<dyn GlobalPanel>) {
        self.right = Some(panel);
    }

    pub fn active(&self) -> Option<&Box<dyn GlobalPanel>> {
        match self.active {
            0 => if let Some(bottom) = &self.bottom {
                Some(&bottom)
            } else {
                None
            },
            1 => if let Some(left) = &self.left {
                Some(&left)
            } else {
                None
            },
            2 => if let Some(right) = &self.right {
                Some(&right)
            } else {
                None
            },
            _ => None,
        }
    }

    pub fn render(&self, child_workspace: Option<&Workspace>, frame: &mut Frame, area: Rect) {
        render(self, child_workspace, frame, area);
    }
}