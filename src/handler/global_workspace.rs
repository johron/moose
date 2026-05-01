use ratatui::Frame;
use ratatui::layout::Rect;
use crate::handler::workspace::Workspace;
use crate::panel::global_panel::GlobalPanel;
use crate::panel::panel::Panel;

#[derive(Debug)]
pub struct GlobalWorkspace {
    panels: Vec<Box<dyn GlobalPanel>>,
    active: usize,
}

impl GlobalWorkspace {
    pub fn new() -> Self {
        GlobalWorkspace {
            panels: Vec::new(),
            active: 0,
        }
    }

    pub fn add_panel(&mut self, panel: Box<dyn GlobalPanel>, make_active: bool) {
        self.panels.push(panel);
        if make_active {
            self.active = self.panels.len() - 1;
        }
        self.panels.last_mut().unwrap().init();
    }

    pub fn active_panel(&self) -> Option<&Box<dyn GlobalPanel>> {
        self.panels.get(self.active)
    }

    pub fn render(&self, child_workspace: Option<&Workspace>, frame: &mut Frame, area: Rect) {
        if let Some(workspace) = child_workspace {
            workspace.render(frame, area);
        }

        for panel in self.panels.iter() {
            if panel.is_shown() == false {
                continue;
            }
            
            let width = area.width as f64 * 1.0;
            let height = area.height as f64 * 0.04;

            panel.render(frame, Rect::new(0, area.height - height as u16, width as u16, height as u16));
        }
    }
}