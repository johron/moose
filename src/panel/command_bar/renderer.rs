use crate::panel::command_bar::command_bar::CommandBar;
use ratatui::layout::Rect;
use ratatui::style::Style;
use ratatui::widgets::{Block, Borders, Paragraph};
use ratatui::Frame;

pub fn render(command_bar: &CommandBar, frame: &mut Frame, area: Rect) {
    let block = Block::default()
        .style(Style::default())
        .borders(Borders::TOP)
        .title("COMMAND");

    let paragraph = Paragraph::new(format!("q:{}", command_bar.cmd.as_str()))
        .block(block)
        .style(Style::default().fg(ratatui::style::Color::Gray));

    frame.render_widget(paragraph, area);
}