use crate::handler::global_workspace::global_workspace::GlobalWorkspace;
use crate::handler::workspace::workspace::Workspace;
use ratatui::layout::{Constraint, Direction, Layout, Rect};
use ratatui::Frame;

pub fn render(global_workspace: &GlobalWorkspace, child_workspace: Option<&Workspace>, frame: &mut Frame, area: Rect) {
    let vertical_constraints = if global_workspace.bottom.is_some() {
        vec![Constraint::Min(1), Constraint::Length(3)]
    } else {
        vec![Constraint::Min(1)]
    };

    let vertical = Layout::default()
        .direction(Direction::Vertical)
        .constraints(vertical_constraints)
        .split(area);

    let top = vertical[0];
    let bottom_area = if global_workspace.bottom.is_some() {
        Some(vertical[1])
    } else {
        None
    };

    let mut horizontal_constraints = Vec::new();

    if global_workspace.left.is_some() {
        horizontal_constraints.push(Constraint::Length(20));
    }

    horizontal_constraints.push(Constraint::Min(1)); // main always exists

    if global_workspace.right.is_some() {
        horizontal_constraints.push(Constraint::Length(20));
    }

    let horizontal = Layout::default()
        .direction(Direction::Horizontal)
        .constraints(horizontal_constraints)
        .split(top);

    let mut i = 0;

    let left_area = if global_workspace.left.is_some() {
        let a = horizontal[i];
        i += 1;
        Some(a)
    } else {
        None
    };

    let main_area = {
        let a = horizontal[i];
        i += 1;
        a
    };

    let right_area = if global_workspace.right.is_some() {
        Some(horizontal[i])
    } else {
        None
    };

    if let Some(workspace) = child_workspace {
        workspace.render(frame, main_area);
    }

    if let (Some(panel), Some(area)) = (&global_workspace.left, left_area) {
        panel.render(frame, area);
    }

    if let (Some(panel), Some(area)) = (&global_workspace.right, right_area) {
        panel.render(frame, area);
    }

    if let (Some(panel), Some(area)) = (&global_workspace.bottom, bottom_area) {
        panel.render(frame, area);
    }
}
