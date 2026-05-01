use ratatui::Frame;
use ratatui::layout::Rect;
use crate::handler::workspace::Workspace;
use crate::panel::global_panel::GlobalPanel;

#[derive(Debug)]
pub struct GlobalWorkspace {
    panels: Vec<Box<dyn GlobalPanel>>,
    shown_vec: Vec<bool>,
    active: usize,
}

impl GlobalWorkspace {
    pub fn new() -> Self {
        GlobalWorkspace {
            panels: Vec::new(),
            shown_vec: Vec::new(),
            active: 0,
        }
    }

    pub fn add_panel(&mut self, panel: Box<dyn GlobalPanel>, show: bool) {
        self.panels.push(panel);
        self.shown_vec.push(show);
    }

    pub fn render(&self, child_workspace: Option<&Workspace>, frame: &mut Frame, area: Rect) {
        if let Some(workspace) = child_workspace {
            workspace.render(frame, area);
        }

        for panel in self.panels.iter() {
            let width = area.width as f64 * 1.0;
            let height = area.height as f64 * 0.04;

            panel.render(frame, Rect::new(0, area.height - height as u16, width as u16, height as u16));
        }
    }
}