use config::Config;
use serde::{Deserialize, Serialize};
use crate::handler::config::{ensure_config_exists, make_config};

pub fn init_config() -> Result<EditorConfig, config::ConfigError> {
    ensure_config_exists::<EditorConfig>(String::from("editor.toml")).expect("Failed to ensure config file exists");
    let cfg = make_config(String::from("editor.toml"), true)?;

    let editor: EditorConfig = cfg.try_deserialize()?;
    Ok(editor)
}

#[derive(Debug, Deserialize, Serialize, Clone)]
pub struct EditorConfig {
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