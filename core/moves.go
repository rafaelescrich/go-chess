package core

import (
	"fmt"
	"strings"
)

/* We could have one struct with enough information to represent the different
   types of moves, but it felt cleaner to use multiple types. */
type move interface {
	perform_move(board *Board) *Board
	String() string
}

type xy struct {
	x, y int
}

type abstract_move struct {
	from, to xy
}

type normal_move struct {
	abstract_move
}

/* Used for keeping track of en-passant */
type pawn_push_two struct {
	abstract_move
}

type pawn_promote struct {
	abstract_move
	piece piece
}

type pawn_en_passant struct {
	abstract_move
}

type king_castle struct {
	king abstract_move
	rook abstract_move
}

type moves []move

func (m normal_move) perform_move(board *Board) *Board {
	b := *board // make a copy of the board
	b.board[m.to.x][m.to.y] = b.board[m.from.x][m.from.y]
	b.board[m.from.x][m.from.y] = square{nil, White}
	b.en_passant = -2
	b.to_play = flip_color(b.to_play)
	return &b
}

func (m pawn_push_two) perform_move(board *Board) *Board {
	b := *board // make a copy of the board
	b.board[m.to.x][m.to.y] = b.board[m.from.x][m.from.y]
	b.board[m.from.x][m.from.y] = square{nil, White}
	assert(m.to.x == m.from.x, "pawn push_push_to not in straight line")
	b.en_passant = m.to.x
	b.to_play = flip_color(b.to_play)
	return &b
}

func (m pawn_promote) perform_move(board *Board) *Board {
	b := *board // make a copy of the board
	b.board[m.to.x][m.to.y] = b.board[m.from.x][m.from.y]
	b.board[m.to.x][m.to.y].piece = m.piece
	b.board[m.from.x][m.from.y] = square{nil, White}
	b.en_passant = -2
	b.to_play = flip_color(b.to_play)
	return &b
}

func (m pawn_en_passant) perform_move(board *Board) *Board {
	b := *board // make a copy of the board
	b.board[m.to.x][m.to.y] = b.board[m.from.x][m.from.y]
	b.board[m.from.x][m.from.y] = square{nil, White}
	b.board[m.to.x][m.from.y] = square{nil, White}
	b.en_passant = -2
	b.to_play = flip_color(b.to_play)
	return &b
}

func is_empty(board *Board, point xy) bool {
	return board.board[point.x][point.y].piece == nil
}

func is_opponent(board *Board, point xy) bool {
	return board.board[point.x][point.y].piece != nil &&
		board.board[point.x][point.y].color != board.to_play
}

func append_if_not_in_check(board *Board, moves moves, move move) (moves, bool) {
	// TODO: perform check analysis
	return append(moves, move), true
}

/**
 * Handles piece movement for all pieces except pawns
 */
func list_moves_common(board *Board, point xy, offsets []xy, repeat bool) moves {
	r := make(moves, 0, 30)
	for _, offset := range offsets {
		r = list_moves_direction(r, board, point, offset, repeat)
	}
	return r
}

func list_moves_direction(r moves, board *Board, point xy, offset xy, repeat bool) moves {
	for i := 1; i < 8; i++ {
		to := xy{point.x + offset.x*i, point.y + offset.y*i}
		if to.x < 0 || to.x >= 8 || to.y < 0 || to.y >= 8 {
			// we hit the edge of the board
			return r
		}
		if board.board[to.x][to.y].piece == nil ||
			board.board[to.x][to.y].color != board.to_play {
			// we are moving into an empty square or capturing a piece
			var ok bool
			r, ok = append_if_not_in_check(board, r,
				normal_move{abstract_move{point, to}})
			if !ok {
				// moving this piece caused a check, so we are done.
				return r
			}
		}
		if !repeat || board.board[to.x][to.y].piece != nil {
			// we are done
			return r
		}
	}
	return r
}

func flip_color(c color) color {
	if c == White {
		return Black
	} else {
		return White
	}
}

func (m normal_move) String() string {
	return fmt.Sprintf("%c%d→%c%d", 'a'+m.from.x, m.from.y+1, 'a'+m.to.x,
		m.to.y+1)
}

func (m pawn_push_two) String() string {
	return fmt.Sprintf("%c%d→%c%d", 'a'+m.from.x, m.from.y+1, 'a'+m.to.x,
		m.to.y+1)
}

func (m pawn_promote) String() string {
	return fmt.Sprintf("%c%d→%c%d%s", 'a'+m.from.x, m.from.y+1, 'a'+m.to.x,
		m.to.y+1, m.piece.String())
}

func (m pawn_en_passant) String() string {
	return fmt.Sprintf("%c%d→%c%de.p.", 'a'+m.from.x, m.from.y+1, 'a'+m.to.x,
		m.to.y+1)
}

func (moves moves) String() string {
	r := make([]string, len(moves))
	for k, m := range moves {
		r[k] = m.String()
	}
	return strings.Join(r, ",")
}
