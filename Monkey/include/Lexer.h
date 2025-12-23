#ifndef RECURSIVEDESCENT_LEXER_H
#define RECURSIVEDESCENT_LEXER_H

#include "Token.h"
#include <cctype>
#include <iosfwd>
#include <istream>
#include <string>

class Lexer
{
public:
	explicit Lexer(std::istream& input)
		: m_input(input)
		, m_peekedToken()
		, m_hasPeeked(false)
	{
	}

	Token GetNextToken()
	{
		if (m_hasPeeked)
		{
			m_hasPeeked = false;
			return m_peekedToken;
		}
		return ReadToken();
	}

	Token PeekToken()
	{
		if (!m_hasPeeked)
		{
			m_peekedToken = ReadToken();
			m_hasPeeked = true;
		}
		return m_peekedToken;
	}

private:
	std::istream& m_input;
	Token m_peekedToken;
	bool m_hasPeeked;

	[[nodiscard]] Token ReadToken() const
	{
		char ch;
		while (m_input.get(ch) && isspace(ch))
		{
		}

		if (m_input.eof())
		{
			return { TokenType::T_EOF, "" };
		}

		if (isalpha(ch))
		{
			return ReadIdentifier(ch);
		}
		if (isdigit(ch))
		{
			return ReadNumber(ch);
		}

		return ReadSymbol(ch);
	}

	[[nodiscard]] Token ReadIdentifier(const char firstChar) const
	{
		std::string s(1, firstChar);
		char ch;
		while (m_input.get(ch) && isalnum(ch))
		{
			s += ch;
		}
		m_input.putback(ch);
		return { GetKeywordType(s), s };
	}

	[[nodiscard]] static TokenType GetKeywordType(const std::string& s)
	{
		if (s == "main")
		{
			return TokenType::T_MAIN;
		}
		if (s == "end")
		{
			return TokenType::T_END;
		}
		if (s == "begin")
		{
			return TokenType::T_BEGIN;
		}
		if (s == "var")
		{
			return TokenType::T_VAR;
		}
		if (s == "int")
		{
			return TokenType::T_INT;
		}
		if (s == "float")
		{
			return TokenType::T_FLOAT;
		}
		if (s == "const")
		{
			return TokenType::T_CONST;
		}
		return TokenType::T_ID;
	}

	[[nodiscard]] Token ReadNumber(const char firstChar) const
	{
		std::string s(1, firstChar);
		char ch;
		while (m_input.get(ch) && (isdigit(ch) || ch == '.'))
		{
			s += ch;
		}
		m_input.putback(ch);
		return { TokenType::T_NUMBER, s };
	}

	[[nodiscard]] Token ReadSymbol(const char ch) const
	{
		if (ch == ':')
			return ReadColon();

		switch (ch)
		{
		case ';':
			return { TokenType::T_SEMICOLON, ";" };
		case ',':
			return { TokenType::T_COMMA, "," };
		case '.':
			return { TokenType::T_DOT, "." };
		case '=':
			return { TokenType::T_EQUALS, "=" };
		case '+':
			return { TokenType::T_PLUS, "+" };
		case '*':
			return { TokenType::T_MULTIPLICATION, "*" };
		case '-':
			return { TokenType::T_MINUS, "-" };
		case '(':
			return { TokenType::T_LEFT_PARENTHESIS, "(" };
		case ')':
			return { TokenType::T_RIGHT_PARENTHESIS, ")" };
		default:
			return { TokenType::T_ERROR, std::string(1, ch) };
		}
	}

	[[nodiscard]] Token ReadColon() const
	{
		char ch;
		if (m_input.get(ch) && ch == '=')
		{
			return { TokenType::T_ASSIGN, ":=" };
		}
		m_input.putback(ch);
		return { TokenType::T_COLON, ":" };
	}
};

#endif // RECURSIVEDESCENT_LEXER_H