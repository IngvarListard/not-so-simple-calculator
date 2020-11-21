package calc

import "testing"

func TestSolver_Expr(t *testing.T) {
	type fields struct {
		parser        *Parser
		currentLexeme Lexeme
	}
	tests := []struct {
		name    string
		fields  fields
		want    interface{}
		wantErr bool
	}{
		{
			name: "simple_sum",
			fields: fields{
				parser: &Parser{
					text:        []rune("3 + 2"),
					pos:         0,
					currentRune: '3',
				},
				currentLexeme: nil,
			},
			want:    int64(5),
			wantErr: false,
		},
		{
			name: "simple_sub",
			fields: fields{
				parser: &Parser{
					text:        []rune("3 - 2"),
					pos:         0,
					currentRune: '3',
				},
				currentLexeme: nil,
			},
			want:    int64(1),
			wantErr: false,
		},
		{
			name: "simple_mul",
			fields: fields{
				parser: &Parser{
					text:        []rune("3 * 2"),
					pos:         0,
					currentRune: '3',
				},
				currentLexeme: nil,
			},
			want:    int64(6),
			wantErr: false,
		},
		{
			name: "simple_div",
			fields: fields{
				parser: &Parser{
					text:        []rune("6 / 3"),
					pos:         0,
					currentRune: '6',
				},
				currentLexeme: nil,
			},
			want:    float64(2),
			wantErr: false,
		},
		{
			name: "complex",
			fields: fields{
				parser: &Parser{
					text:        []rune(" 3 + 2 * 6 / 4 "),
					pos:         0,
					currentRune: ' ',
				},
				currentLexeme: nil,
			},
			want:    float64(6),
			wantErr: false,
		},
		{
			name: "simple_brackets",
			fields: fields{
				parser: &Parser{
					text:        []rune(" 2 * (2 + 5) "),
					pos:         0,
					currentRune: ' ',
				},
				currentLexeme: nil,
			},
			want:    int64(14),
			wantErr: false,
		},
		{
			name: "brackets_formula",
			fields: fields{
				parser: &Parser{
					text:        []rune("2 * ((2 + 5) + ((3 * 8) * (2 + 5))"),
					pos:         0,
					currentRune: '2',
				},
				currentLexeme: nil,
			},
			want:    int64(350),
			wantErr: false,
		},
		{
			name: "brackets_formula_3",
			fields: fields{
				parser: &Parser{
					text:        []rune("(8 + 10) / 2"),
					pos:         0,
					currentRune: '(',
				},
				currentLexeme: nil,
			},
			want:    float64(9),
			wantErr: false,
		},
		{
			name: "brackets_formula_2",
			fields: fields{
				parser: &Parser{
					text:        []rune("7 + 3 * (10 / (12 / (3 + 1) - 1))"),
					pos:         0,
					currentRune: '7',
				},
				currentLexeme: nil,
			},
			want:    float64(22),
			wantErr: false,
		},
		{
			name: "parsing_error_1",
			fields: fields{
				parser: &Parser{
					text:        []rune("2 5"),
					pos:         0,
					currentRune: '2',
				},
				currentLexeme: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "simple_float_sum",
			fields: fields{
				parser: &Parser{
					text:        []rune("2.2 + 5"),
					pos:         0,
					currentRune: '2',
				},
				currentLexeme: nil,
			},
			want:    7.2,
			wantErr: false,
		},
		{
			name: "simple_float_mul",
			fields: fields{
				parser: &Parser{
					text:        []rune("2.2 * 2"),
					pos:         0,
					currentRune: '2',
				},
				currentLexeme: nil,
			},
			want:    4.4,
			wantErr: false,
		},
		{
			name: "simple_float_div",
			fields: fields{
				parser: &Parser{
					text:        []rune("4.4 / 2"),
					pos:         0,
					currentRune: '4',
				},
				currentLexeme: nil,
			},
			want:    2.2,
			wantErr: false,
		},
		{
			name: "simple_float_sub",
			fields: fields{
				parser: &Parser{
					text:        []rune("4.0 - 2"),
					pos:         0,
					currentRune: '4',
				},
				currentLexeme: nil,
			},
			want:    2.0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Solver{
				parser:        tt.fields.parser,
				currentLexeme: tt.fields.currentLexeme,
			}
			var err error
			i.parser.skipWhitespace()
			i.currentLexeme, err = i.parser.getNextLexeme()
			if err != nil {
				t.Fatalf("%v", err)
				return
			}
			got, err := i.Solve()
			if i.parser.currentRune != -1 {
				t.Errorf("NOT EOF")
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Expr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Expr() got = %v, want %v", got, tt.want)
			}
		})
	}
}
