use crate::handler::config::{ensure_config_exists, make_config};
use serde::{Deserialize, Serialize};

pub fn init_config() -> Result<EditorConfig, config::ConfigError> {
    ensure_config_exists::<EditorConfig>(String::from("builtin/editor.toml")).expect("Failed to ensure config file exists");
    let cfg = make_config(String::from("builtin/editor.toml"), true)?;

    let editor: EditorConfig = cfg.try_deserialize()?;
    Ok(editor)
}

#[derive(Debug, Deserialize, Serialize, Clone)]
pub struct EditorConfig {
    pub enable_mouse: bool,
    pub enable_blinking: bool,
    pub blink_interval: f64,
    
    pub enter_normal_mode: String,
    pub enter_insert_mode: String,
    pub enter_command_mode: String,

    pub insert_mode_delete_char: String,
    pub insert_mode_newline: String,
    pub insert_mode_indent: String,
    pub insert_mode_reverse_indent: String,
}

impl Default for EditorConfig {
    fn default() -> Self {
        EditorConfig {
            enable_mouse: true,
            enable_blinking: false,
            blink_interval: 0.5,
            
            enter_normal_mode: String::from("esc"),
            enter_insert_mode: String::from("char:i"),
            enter_command_mode: String::from("char:q"),

            insert_mode_delete_char: String::from("backspace"),
            insert_mode_newline: String::from("enter"),
            insert_mode_indent: String::from("tab"),
            insert_mode_reverse_indent: String::from("shift+tab"),
        }
    }
}