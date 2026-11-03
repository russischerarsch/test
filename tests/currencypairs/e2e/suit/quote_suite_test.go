package suit

import (
	"PlataTest/tests/currencypairs/e2e/checks"
	"PlataTest/tests/currencypairs/models"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

func (s *QuoteSuite) TestCreateQuote_success() {
	t := s.T()
	resp, err := s.Client.CreateQuote(t.Context(), models.CreatePairReq{Pair: "EURUSD"})
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal(http.StatusCreated, resp.Code)
	t.Cleanup(func() {
		s.DeleteQuote(resp.Data.ID)
	})
	checks.CreatePairResponseValid(t, resp.Data)
	quote, err := s.Client.GetQuoteByID(t.Context(), resp.Data.ID)
	s.Require().NoError(err)
	s.Require().NotNil(quote)
	s.Assert().Equal("EURUSD", quote.Data.Pair)
}

func TestRunSuit(t *testing.T) {
	suite.Run(t, new(QuoteSuite))
}
