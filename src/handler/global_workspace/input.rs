use crate::handler::config::{char_config, vec_config};
use crate::handler::global_workspace::global_workspace::{GlobalWorkspace, GlobalWorkspaceActive};
use crate::handler::input::InputEvent;
use crate::handler::workspace::workspace::Workspace;
use crate::panel::global_panel::GlobalPanelMeta;

impl GlobalWorkspace {
    pub fn input(&mut self, child_workspace: Option<&mut Workspace>, input_event: InputEvent) {
        // First do global input, if nothing matches then continue

        let res = Self::handle_global_input(self, child_workspace, input_event.clone());
        if res.is_none() {
            return;
        }
        let child_workspace = res.unwrap();

        match self.active {
            GlobalWorkspaceActive::Workspace => {
                if let Some(workspace) = child_workspace {
                    if let Some(panel) = workspace.active_panel_mut() {
                        if panel.get_global_panels().is_some() && panel.is_normal_mode() {
                            let metas = panel.get_global_panels().unwrap();
                            if self.handle_panel_global_panel_input(metas, input_event.clone()) {
                                return;
                            }
                        }

                        panel.input(input_event);
                    }
                }
            },
            GlobalWorkspaceActive::Bottom => if let Some(bottom) = &mut self.bottom {
                bottom.input(input_event);
            },
            GlobalWorkspaceActive::Left => if let Some(left) = &mut self.left {
                left.input(input_event);
            },
            GlobalWorkspaceActive::Right => if let Some(right) = &mut self.right {
                right.input(input_event);
            }
        }
    }

    fn handle_panel_global_panel_input(&mut self, metas: Vec<GlobalPanelMeta>, input_event: InputEvent) -> bool {
        for meta in metas {
            if let InputEvent::Char(c) = input_event.clone() {
                if char_config(meta.toggle_show_shortcut) == Some(c) {
                    match meta.location {
                        GlobalWorkspaceActive::Bottom => if let Some(bottom) = &mut self.bottom {
                            if bottom.identity() == meta.panel.identity() {
                                self.bottom = None;
                                self.active = GlobalWorkspaceActive::Workspace;
                            } else {
                                self.set_bottom(meta.panel);
                                self.active = GlobalWorkspaceActive::Bottom;
                            }
                        } else {
                            self.set_bottom(meta.panel);
                            self.active = GlobalWorkspaceActive::Bottom;

                        },
                        GlobalWorkspaceActive::Left => if let Some(left) = &mut self.left {
                            if left.identity() == meta.panel.identity() {
                                self.left = None;
                                self.active = GlobalWorkspaceActive::Workspace;
                            } else {
                                self.set_left(meta.panel);
                                self.active = GlobalWorkspaceActive::Left;
                            }
                        } else {
                            self.set_left(meta.panel);
                            self.active = GlobalWorkspaceActive::Left;
                        },
                        GlobalWorkspaceActive::Right => if let Some(right) = &mut self.right {
                            if right.identity() == meta.panel.identity() {
                                self.right = None;
                                self.active = GlobalWorkspaceActive::Workspace;
                            } else {
                                self.set_right(meta.panel);
                                self.active = GlobalWorkspaceActive::Right;
                            }
                        } else {
                            self.set_right(meta.panel);
                            self.active = GlobalWorkspaceActive::Right;
                        },
                        _ => {}
                    }

                    return true;
                }
            } else if let InputEvent::Keyboard(vec) = input_event.clone() {
                if vec_config(meta.toggle_show_shortcut) == vec {
                    match meta.location {
                        GlobalWorkspaceActive::Bottom => if let Some(bottom) = &mut self.bottom {
                            if bottom.identity() == meta.panel.identity() {
                                self.bottom = None;
                                self.active = GlobalWorkspaceActive::Workspace;
                            } else {
                                self.set_bottom(meta.panel);
                                self.active = GlobalWorkspaceActive::Bottom;
                            }
                        } else {
                            self.set_bottom(meta.panel);
                            self.active = GlobalWorkspaceActive::Bottom;
                        },
                        GlobalWorkspaceActive::Left => if let Some(left) = &mut self.left {
                            if left.identity() == meta.panel.identity() {
                                self.left = None;
                                self.active = GlobalWorkspaceActive::Workspace;
                            } else {
                                self.set_left(meta.panel);
                                self.active = GlobalWorkspaceActive::Left;
                            }
                        } else {
                            self.set_left(meta.panel);
                            self.active = GlobalWorkspaceActive::Left;
                        },
                        GlobalWorkspaceActive::Right => if let Some(right) = &mut self.right {
                            if right.identity() == meta.panel.identity() {
                                self.right = None;
                                self.active = GlobalWorkspaceActive::Workspace;
                            } else {
                                self.set_right(meta.panel);
                                self.active = GlobalWorkspaceActive::Right;
                            }
                        } else {
                            self.set_right(meta.panel);
                            self.active = GlobalWorkspaceActive::Right;
                        },
                        _ => {}
                    }

                    return true;
                }
            }
        }

        false
    }

    fn handle_global_input<'a>(&mut self, mut child_workspace: Option<&'a mut Workspace>, input_event: InputEvent) -> Option<Option<&'a mut Workspace>> {
        if let Some(ws) = child_workspace.as_mut() {
            if let Some(panel) = ws.active_panel() {
                if !panel.is_normal_mode() {
                    return Some(child_workspace);
                }
            }
        }

        let config = self.config.clone();

        match input_event {
            InputEvent::Char(c) => {
                if Some(c) == char_config(config.enter_global_command_mode) {
                    todo!("Change to GlobalPanel mode, show command global panel")
                } else if Some(c) == char_config(String::from("char:-")) {
                    self.should_quit = true;
                    return None;
                }
            },
            _ => {}
        }

        Some(child_workspace)
    }
}
