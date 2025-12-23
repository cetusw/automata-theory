#pragma once

#include "Lexer.h"
#include "Token.h"

class SecondPopulationParser {
public:
	bool Parse(Lexer& lexer);

private:
	bool parseRule1(Lexer& lexer);
	bool parseRule2(Lexer& lexer);
	bool parseRule3(Lexer& lexer);
	bool expectToken(Lexer& lexer, TokenType expected);
};

inline bool SecondPopulationParser::expectToken(Lexer& lexer, TokenType expected) {
	Token token = lexer.GetNextToken();

	if (token.m_type == TokenType::T_EOF) {
		return false;
	}

	if (token.m_type != expected) {
		return false;
	}

	return true;
}

inline bool SecondPopulationParser::parseRule1(Lexer& lexer) {
	if (!expectToken(lexer, TokenType::T_OI)) {
		return false;
	}

	if (!parseRule2(lexer)) {
		return false;
	}

	if (!expectToken(lexer, TokenType::T_AI)) {
		return false;
	}

	if (!parseRule3(lexer)) {
		return false;
	}

	return true;
}

inline bool SecondPopulationParser::parseRule2(Lexer& lexer) {
	if (!expectToken(lexer, TokenType::T_NU)) {
		return false;
	}

	Token peek = lexer.PeekToken();
	if (peek.m_type == TokenType::T_EOF) {
		return true;
	}

	return parseRule2(lexer);
}

inline bool SecondPopulationParser::parseRule3(Lexer& lexer) {
	Token peek = lexer.PeekToken();
	if (peek.m_type == TokenType::T_EOF) {
		return false;
	}

	if (peek.m_type == TokenType::T_UH_TI) {
		lexer.GetNextToken();
		return true;
	}

	if (peek.m_type == TokenType::T_HO) {
		if (!expectToken(lexer, TokenType::T_HO)) {
			return false;
		}

		if (!parseRule3(lexer)) {
			return false;
		}

		if (!expectToken(lexer, TokenType::T_HO)) {
			return false;
		}

		return true;
	}

	return false;
}

inline bool SecondPopulationParser::Parse(Lexer& lexer) {
	return parseRule1(lexer);
}
