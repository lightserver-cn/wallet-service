package test

import (
	"fmt"
	"net/http"
	"testing"

	"server/app/model"
	"server/app/request"
	"server/pkg/consts"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestWalletsDeposit(t *testing.T) {
	defer goleak.VerifyNone(
		t,
		goleak.IgnoreTopFunction("net/http.(*Server).Serve"),
		goleak.IgnoreTopFunction("net/http/httptest.(*Server).goServe.func1"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).readLoop"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).writeLoop"),
		goleak.IgnoreTopFunction("internal/poll.runtime_pollWait"),
		goleak.IgnoreTopFunction("internal/poll.(*pollDesc).wait"),
		goleak.IgnoreTopFunction("internal/poll.(*FD).Accept"),
		goleak.IgnoreTopFunction("internal/poll.(*FD).Read"),
		goleak.IgnoreTopFunction("time.Sleep"),
		goleak.IgnoreTopFunction("time.AfterFunc"),
		goleak.IgnoreTopFunction("time.Ticker"),
		goleak.IgnoreTopFunction("runtime.gopark"),
		goleak.IgnoreTopFunction("runtime.forcegchelper"),
		goleak.IgnoreTopFunction("runtime.bgsweep"),
		goleak.IgnoreTopFunction("runtime.bgscavenge"),
	)

	m := NewMockTest().start(t)
	defer m.Teardown()

	t.Run("deposit", func(t *testing.T) {
		var uid int64 = 1

		// balance
		resGetBalance := m.Expect.GET(fmt.Sprintf("/api/wallets/%d/balance", uid)).Expect().Status(http.StatusOK).JSON()

		respGetBalance := &request.ResBalance{}
		if err := AssertResponse(resGetBalance.Raw(), &respGetBalance); err != nil {
			t.Error(err)
		}
		beforeBalance := respGetBalance.Balance

		// deposit
		var amount int64 = 10
		req := map[string]any{
			"amount": amount,
		}
		resDeposit := m.Expect.POST(fmt.Sprintf("/api/wallets/%d/deposit", uid)).WithJSON(req).Expect().Status(http.StatusOK).JSON()
		AssertResponseSuccess(t, consts.MsgSuccess, resDeposit, "message mismatch")

		// balance
		resGetBalance = m.Expect.GET(fmt.Sprintf("/api/wallets/%d/balance", uid)).Expect().Status(http.StatusOK).JSON()

		respGetBalance = &request.ResBalance{}
		if err := AssertResponse(resGetBalance.Raw(), &respGetBalance); err != nil {
			t.Error(err)
		}
		afterBalance := respGetBalance.Balance

		amountActual := afterBalance.Sub(beforeBalance)
		amountExpected := decimal.NewFromInt(amount)

		assert.Equal(t, amountExpected, amountActual, "deposit mismatch")
	})
}

func TestWalletsWithDraw(t *testing.T) {
	defer goleak.VerifyNone(
		t,
		goleak.IgnoreTopFunction("net/http.(*Server).Serve"),
		goleak.IgnoreTopFunction("net/http/httptest.(*Server).goServe.func1"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).readLoop"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).writeLoop"),
		goleak.IgnoreTopFunction("internal/poll.runtime_pollWait"),
		goleak.IgnoreTopFunction("internal/poll.(*pollDesc).wait"),
		goleak.IgnoreTopFunction("internal/poll.(*FD).Accept"),
		goleak.IgnoreTopFunction("internal/poll.(*FD).Read"),
		goleak.IgnoreTopFunction("time.Sleep"),
		goleak.IgnoreTopFunction("time.AfterFunc"),
		goleak.IgnoreTopFunction("time.Ticker"),
		goleak.IgnoreTopFunction("runtime.gopark"),
		goleak.IgnoreTopFunction("runtime.forcegchelper"),
		goleak.IgnoreTopFunction("runtime.bgsweep"),
		goleak.IgnoreTopFunction("runtime.bgscavenge"),
	)

	m := NewMockTest().start(t)
	defer m.Teardown()

	t.Run("withdraw", func(t *testing.T) {
		var uid int64 = 1

		// balance
		resGetBalance := m.Expect.GET(fmt.Sprintf("/api/wallets/%d/balance", uid)).Expect().Status(http.StatusOK).JSON()
		respGetBalance := &request.ResBalance{}
		if err := AssertResponse(resGetBalance.Raw(), &respGetBalance); err != nil {
			t.Error(err)
		}
		beforeBalance := respGetBalance.Balance

		// withdraw
		var amount int64 = 10
		req := map[string]any{
			"amount": amount,
		}
		resWithdraw := m.Expect.POST(fmt.Sprintf("/api/wallets/%d/withdraw", uid)).WithJSON(req).Expect().Status(http.StatusOK).JSON()
		AssertResponseSuccess(t, consts.MsgSuccess, resWithdraw, "message mismatch")

		// balance
		resGetBalance = m.Expect.GET(fmt.Sprintf("/api/wallets/%d/balance", uid)).Expect().Status(http.StatusOK).JSON()
		respGetBalance = &request.ResBalance{}
		if err := AssertResponse(resGetBalance.Raw(), &respGetBalance); err != nil {
			t.Error(err)
		}
		afterBalance := respGetBalance.Balance

		amountActual := beforeBalance.Sub(afterBalance)
		amountExpected := decimal.NewFromInt(amount)

		assert.Equal(t, amountExpected, amountActual, "withdraw mismatch")
	})
}

func TestWalletsTransfer(t *testing.T) {
	defer goleak.VerifyNone(
		t,
		goleak.IgnoreTopFunction("net/http.(*Server).Serve"),
		goleak.IgnoreTopFunction("net/http/httptest.(*Server).goServe.func1"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).readLoop"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).writeLoop"),
		goleak.IgnoreTopFunction("internal/poll.runtime_pollWait"),
		goleak.IgnoreTopFunction("internal/poll.(*pollDesc).wait"),
		goleak.IgnoreTopFunction("internal/poll.(*FD).Accept"),
		goleak.IgnoreTopFunction("internal/poll.(*FD).Read"),
		goleak.IgnoreTopFunction("time.Sleep"),
		goleak.IgnoreTopFunction("time.AfterFunc"),
		goleak.IgnoreTopFunction("time.Ticker"),
		goleak.IgnoreTopFunction("runtime.gopark"),
		goleak.IgnoreTopFunction("runtime.forcegchelper"),
		goleak.IgnoreTopFunction("runtime.bgsweep"),
		goleak.IgnoreTopFunction("runtime.bgscavenge"),
	)

	m := NewMockTest().start(t)
	defer m.Teardown()

	t.Run("transfer", func(t *testing.T) {
		var fromUID int64 = 1
		var toUID int64 = 2

		// balance
		resGetBalance := m.Expect.GET(fmt.Sprintf("/api/wallets/%d/balance", fromUID)).Expect().Status(http.StatusOK).JSON()
		respGetBalance := &request.ResBalance{}
		if err := AssertResponse(resGetBalance.Raw(), &respGetBalance); err != nil {
			t.Error(err)
		}
		beforeBalanceFrom := respGetBalance.Balance

		// balance
		resGetBalance = m.Expect.GET(fmt.Sprintf("/api/wallets/%d/balance", toUID)).Expect().Status(http.StatusOK).JSON()
		respGetBalance = &request.ResBalance{}
		if err := AssertResponse(resGetBalance.Raw(), &respGetBalance); err != nil {
			t.Error(err)
		}
		beforeBalanceTo := respGetBalance.Balance

		// transfer
		var amount int64 = 2
		req := map[string]any{
			"to_uid": toUID,
			"amount": amount,
		}
		resTransfer := m.Expect.POST(fmt.Sprintf("/api/wallets/%d/transfer", fromUID)).WithJSON(req).Expect().Status(http.StatusOK).JSON()
		AssertResponseSuccess(t, consts.MsgSuccess, resTransfer, "message mismatch")

		// balance
		resGetBalance = m.Expect.GET(fmt.Sprintf("/api/wallets/%d/balance", fromUID)).Expect().Status(http.StatusOK).JSON()
		respGetBalance = &request.ResBalance{}
		if err := AssertResponse(resGetBalance.Raw(), &respGetBalance); err != nil {
			t.Error(err)
		}
		afterBalanceFrom := respGetBalance.Balance

		// balance
		resGetBalance = m.Expect.GET(fmt.Sprintf("/api/wallets/%d/balance", toUID)).Expect().Status(http.StatusOK).JSON()
		respGetBalance = &request.ResBalance{}
		if err := AssertResponse(resGetBalance.Raw(), &respGetBalance); err != nil {
			t.Error(err)
		}
		afterBalanceTo := respGetBalance.Balance

		amountActualFrom := beforeBalanceFrom.Sub(afterBalanceFrom)
		amountActualTo := afterBalanceTo.Sub(beforeBalanceTo)
		amountExpected := decimal.NewFromInt(amount)

		assert.Equal(t, amountExpected, amountActualFrom, "withdraw mismatch")
		assert.Equal(t, amountExpected, amountActualTo, "withdraw mismatch")
	})
}

func TestWalletsBalance(t *testing.T) {
	defer goleak.VerifyNone(
		t,
		goleak.IgnoreTopFunction("net/http.(*Server).Serve"),
		goleak.IgnoreTopFunction("net/http/httptest.(*Server).goServe.func1"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).readLoop"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).writeLoop"),
		goleak.IgnoreTopFunction("internal/poll.runtime_pollWait"),
		goleak.IgnoreTopFunction("internal/poll.(*pollDesc).wait"),
		goleak.IgnoreTopFunction("internal/poll.(*FD).Accept"),
		goleak.IgnoreTopFunction("internal/poll.(*FD).Read"),
		goleak.IgnoreTopFunction("time.Sleep"),
		goleak.IgnoreTopFunction("time.AfterFunc"),
		goleak.IgnoreTopFunction("time.Ticker"),
		goleak.IgnoreTopFunction("runtime.gopark"),
		goleak.IgnoreTopFunction("runtime.forcegchelper"),
		goleak.IgnoreTopFunction("runtime.bgsweep"),
		goleak.IgnoreTopFunction("runtime.bgscavenge"),
	)

	m := NewMockTest().start(t)
	defer m.Teardown()

	t.Run("balance", func(t *testing.T) {
		var uid int64 = 1

		// balance
		resGetBalance := m.Expect.GET(fmt.Sprintf("/api/wallets/%d/balance", uid)).Expect().Status(http.StatusOK).JSON()

		respGetBalance := &request.ResBalance{}
		if err := AssertResponse(resGetBalance.Raw(), &respGetBalance); err != nil {
			t.Error(err)
		}

		assert.Equal(t, decimal.NewFromInt(58), respGetBalance.Balance, "balance mismatch")
	})
}

func TestWalletsTransactions(t *testing.T) {
	defer goleak.VerifyNone(
		t,
		goleak.IgnoreTopFunction("net/http.(*Server).Serve"),
		goleak.IgnoreTopFunction("net/http/httptest.(*Server).goServe.func1"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).readLoop"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).writeLoop"),
		goleak.IgnoreTopFunction("internal/poll.runtime_pollWait"),
		goleak.IgnoreTopFunction("internal/poll.(*pollDesc).wait"),
		goleak.IgnoreTopFunction("internal/poll.(*FD).Accept"),
		goleak.IgnoreTopFunction("internal/poll.(*FD).Read"),
		goleak.IgnoreTopFunction("time.Sleep"),
		goleak.IgnoreTopFunction("time.AfterFunc"),
		goleak.IgnoreTopFunction("time.Ticker"),
		goleak.IgnoreTopFunction("runtime.gopark"),
		goleak.IgnoreTopFunction("runtime.forcegchelper"),
		goleak.IgnoreTopFunction("runtime.bgsweep"),
		goleak.IgnoreTopFunction("runtime.bgscavenge"),
	)

	m := NewMockTest().start(t)
	defer m.Teardown()

	t.Run("transactions", func(t *testing.T) {
		var uid int64 = 1

		// deposit
		var amountDeposit int64 = 100
		reqDeposit := map[string]any{
			"amount": amountDeposit,
		}
		resDeposit := m.Expect.POST(fmt.Sprintf("/api/wallets/%d/deposit", uid)).WithJSON(reqDeposit).Expect().Status(http.StatusOK).JSON()
		AssertResponseSuccess(t, consts.MsgSuccess, resDeposit, "message mismatch")

		// withdraw
		var amountWithdraw int64 = 10
		reqWithdraw := map[string]any{
			"amount": amountWithdraw,
		}
		resWithdraw := m.Expect.POST(fmt.Sprintf("/api/wallets/%d/withdraw", uid)).WithJSON(reqWithdraw).Expect().Status(http.StatusOK).JSON()
		AssertResponseSuccess(t, consts.MsgSuccess, resWithdraw, "message mismatch")

		// transfer
		var toUID int64 = 2
		var amountTransfer int64 = 12
		req := map[string]any{
			"to_uid": toUID,
			"amount": amountTransfer,
		}
		resTransfer := m.Expect.POST(fmt.Sprintf("/api/wallets/%d/transfer", uid)).WithJSON(req).Expect().Status(http.StatusOK).JSON()
		AssertResponseSuccess(t, consts.MsgSuccess, resTransfer, "message mismatch")

		// transactions
		reqTransaction := map[string]any{
			"page":      1,
			"page_size": 3,
			"type":      0,
		}
		resTransaction := m.Expect.GET(fmt.Sprintf("/api/wallets/%d/transactions", uid)).WithQueryObject(reqTransaction).Expect().Status(http.StatusOK).JSON()
		resp := &request.ResTransactions{
			List:    make([]*model.TransactionWithUsername, 0, 3),
			HasMore: false,
		}
		if err := AssertResponse(resTransaction.Raw(), &resp); err != nil {
			t.Error(err)
		}

		// transfer
		assert.Equal(t, decimal.NewFromInt(amountTransfer), resp.List[0].Amount, "amount mismatch")
		assert.Equal(t, "transfer", resp.List[0].TransactionTypeName, "transaction_type mismatch")

		// withdraw
		assert.Equal(t, decimal.NewFromInt(amountWithdraw), resp.List[1].Amount, "amount mismatch")
		assert.Equal(t, "withdraw", resp.List[1].TransactionTypeName, "transaction_type mismatch")

		// deposit
		assert.Equal(t, decimal.NewFromInt(amountDeposit), resp.List[2].Amount, "amount mismatch")
		assert.Equal(t, "deposit", resp.List[2].TransactionTypeName, "transaction_type mismatch")
	})
}
