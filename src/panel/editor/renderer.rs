use ratatui::Frame;
use ratatui::layout::{Constraint, Direction, Layout, Rect};
use ratatui::style::{Color, Modifier, Style};
use ratatui::text::{Line, Span};
use ratatui::widgets::Paragraph;
use ropey::Rope;
use crate::panel::editor::editor::Editor;

fn make_line_spans(top_line: usize, max_line: usize, rope: &Rope) -> Vec<Line<'static>> {
    let mut lines_spans: Vec<Line> = Vec::new();

    for i in top_line..max_line {
        let rope_line = rope.line(i);
        let text_line = rope_line.to_string();
        let spans = vec![Span::raw(text_line)];
        let mut line_spans = vec![Span::styled(format!("{:4} ", i), Style::default().fg(Color::Gray))]; // small gutter

        line_spans.extend(spans);
        lines_spans.push(Line::from(line_spans));
    }

    lines_spans
}

pub fn render(editor: &Editor, frame: &mut Frame, area: Rect) {
    let chunks = Layout::default()
        .direction(Direction::Horizontal)
        .constraints([Constraint::Percentage(100)].as_ref())
        .split(area);

    let height = std::cmp::max(1, chunks[0].height as usize - 1);

    let max_line = std::cmp::min(
        editor.get_rope().len_lines(),
        editor.get_scroll_offset().saturating_add(height),
    );

    let lines_spans: Vec<Line> = make_line_spans(editor.get_scroll_offset(), max_line, editor.get_rope());

    // Have to think about how I can to the multiple editor panels later, block should be set from outside, not in editor panel
    let paragraph = Paragraph::new(lines_spans);
    //.bg(Color::Red);
    //.block(block);

    frame.render_widget(paragraph, chunks[0]);

    for cursor in editor.get_cursors().iter() {
        let x = chunks[0].x + 5 + cursor.col as u16; // 5 for gutter
        let y = chunks[0].y + (cursor.line.saturating_sub(editor.get_scroll_offset())) as u16;
        frame.render_widget(
            Paragraph::new("")
                .style(Style::default()
                    .add_modifier(Modifier::REVERSED)
                    .fg(Color::White)),
            //.fg(if panel.active { Color::White } else { Color::DarkGray })),
            Rect::new(x, y, 1, 1),
        );
    }
}