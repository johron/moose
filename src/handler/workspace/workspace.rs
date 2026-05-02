use crate::panel::panel::Panel;
use ratatui::layout::Rect;
use ratatui::Frame;

#[derive(Debug)]
pub struct Workspace {
    panels: Vec<Box<dyn Panel>>,
    active: usize,
    direction: Direction,
}

#[derive(Debug, Clone)]
enum Direction {
    Horizontal,
    Vertical,
}

impl Workspace {
    pub fn new() -> Self {
        Workspace {
            panels: Vec::new(),
            active: 0,
            direction: Direction::Horizontal,
        }
    }

    pub fn add_panel(&mut self, panel: Box<dyn Panel>, make_active: bool) {
        self.panels.push(panel);
        if make_active {
            self.active = self.panels.len() - 1;
        }
        self.panels.last_mut().unwrap().init();
    }

    pub fn active_panel_mut(&mut self) -> Option<&mut Box<dyn Panel>> {
        self.panels.get_mut(self.active)
    }
    
    pub fn active_panel(&self) -> Option<&Box<dyn Panel>> {
        self.panels.get(self.active)
    }

    pub fn render(&self, frame: &mut Frame, area: Rect) {
        // Add layout stuff, direction, multiple panels on one workspace
        if let Some(panel) = self.active_panel() {
            panel.render(frame, area);
        }
    }
}