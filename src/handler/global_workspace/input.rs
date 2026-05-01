use crate::handler::config::char_config;
use crate::handler::global_workspace::global_workspace::{GlobalWorkspace, GlobalWorkspaceActive};
use crate::handler::input::InputEvent;
use crate::handler::workspace::workspace::Workspace;

pub fn input(global_workspace: &mut GlobalWorkspace, child_workspace: Option<&mut Workspace>, input_event: InputEvent) {
    // First do global input, if nothing matches then continue

    let res = handle_global_input(global_workspace, child_workspace, input_event.clone());
    if res.is_none() {
        return;
    }
    let (global_workspace, child_workspace) = res.unwrap();

    match global_workspace.active {
        GlobalWorkspaceActive::Workspace => {
            if let Some(workspace) = child_workspace {
                if let Some(panel) = workspace.active_panel_mut() {
                    panel.input(input_event);
                }
            }

            /*let config = global_workspace.config.clone();

            if let Some(workspace) = child_workspace {
                if let Some(panel) = workspace.active_panel_mut() {
                    if panel.is_normal_mode() {
                        match input_event {
                            InputEvent::Char(c) => {
                                if Some(c) == char_config(config.enter_command_mode) {
                                    todo!("Change to GlobalPanel mode, show command global panel")
                                } else if Some(c) == char_config(String::from("char:-")) {
                                    global_workspace.should_quit = true;
                                } else {
                                    panel.input(input_event);
                                }
                            },
                            _ => panel.input(input_event)

                        }
                    } else {
                        panel.input(input_event);
                    }
                }
            }*/
        },
        GlobalWorkspaceActive::Bottom => if let Some(bottom) = &mut global_workspace.bottom {
            bottom.input(input_event);
        },
        GlobalWorkspaceActive::Left => if let Some(left) = &mut global_workspace.left {
            left.input(input_event);
        },
        GlobalWorkspaceActive::Right => if let Some(right) = &mut global_workspace.right {
            right.input(input_event);
        }
    }
}

fn handle_global_input<'a>(global_workspace: &'a mut GlobalWorkspace, child_workspace: Option<&'a mut Workspace>, input_event: InputEvent) -> Option<(&'a mut GlobalWorkspace, Option<&'a mut Workspace>)> {
    let config = global_workspace.config.clone();

    match input_event {
        InputEvent::Char(c) => {
            if Some(c) == char_config(config.enter_command_mode) {
                todo!("Change to GlobalPanel mode, show command global panel")
            } else if Some(c) == char_config(String::from("char:-")) {
                global_workspace.should_quit = true;
                return None;
            }
        },
        _ => {}
    }

    Some((global_workspace, child_workspace))
}