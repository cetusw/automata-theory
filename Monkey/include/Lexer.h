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

		return ReadIdentifier(ch);
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
		if (s == "ну")
		{
			std::cout << "ну" << std::endl;
		}
		if (s == "ау")
		{
			return TokenType::T_AU;
		}
		if (s == "ку")
		{
			return TokenType::T_KU;
		}
		if (s == "ух-ты")
		{
			return TokenType::T_UH_TI;
		}
		if (s == "хо")
		{
			return TokenType::T_HO;
		}
		if (s == "ну")
		{
			return TokenType::T_NU;
		}
		if (s == "и_ну")
		{
			return TokenType::T_I_NU;
		}
		if (s == "ой")
		{
			return TokenType::T_OI;
		}
		if (s == "ай")
		{
			return TokenType::T_AI;
		}
		return TokenType::T_EOF;
	}
};

#endif // RECURSIVEDESCENT_LEXER_H