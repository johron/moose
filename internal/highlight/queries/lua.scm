;; ============================================================================
;; Tree-sitter Query File for Lua (highlights.scm)
;; ============================================================================

;; Keywords
[
  "and" "do" "else" "elseif" "end" 
  "for" "function" "goto" "if" "in" 
  "local" "not" "or" "repeat" "return" 
  "then" "until" "while"
] @keyword

;; Literals and Constants
(string) @string
(number) @number

;; Comments
[
  (comment)
  (hash_bang_line)
] @comment

;; Parameters and Variables
(parameters (identifier) @variable.parameter)

;; Table Fields
(field name: (identifier) @property)
(dot_index_expression field: (identifier) @property)

;; Punctuation & Operators
[ "," "." ";" ":" "#" ] @punctuation.delimiter
[ "(" ")" "{" "}" "[" "]" ] @punctuation.bracket
[ "=" "==" "~=" "<" "<=" ">" ">=" "+" "-" "*" "/" "//" "^" "%" ".." ] @operator
