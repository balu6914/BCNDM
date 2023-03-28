package subscriptions

import (
	"strconv"
	"time"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

const (
	income  = "income"
	outcome = "outcome"
)

var (
	whiteColor = color.NewWhite()

	headerBg = color.Color{
		Red:   6,
		Green: 210,
		Blue:  216,
	}

	incomeBf = color.Color{
		Red:   45,
		Green: 220,
		Blue:  50,
	}

	outcomeBg = color.Color{
		Red:   220,
		Green: 45,
		Blue:  50,
	}

	grayColor = color.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
)

func (ss subscriptionsService) Report(q Query, owner string) ([]byte, error) {
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Report", props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
	})
	m.SetBackgroundColor(grayColor)

	m.Row(9, func() {
		m.Col(2, func() {
			m.Text("Stream name", props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
		m.Col(1, func() {
			m.Text("Type", props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
		m.Col(4, func() {
			m.Text("Start", props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
		m.Col(4, func() {
			m.Text("End", props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
		m.Col(1, func() {
			m.Text("Price", props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
	})

	var totalIncome, totalOutcome uint64
	subs, err := ss.subscriptions.Search(q)
	if err != nil {
		return nil, err
	}

	m.SetBackgroundColor(whiteColor)

	for _, sub := range subs.Content {
		t := income
		if sub.StreamOwner != owner {
			t = outcome
		}
		switch t {
		case income:
			totalIncome += sub.StreamPrice * sub.Hours
		default:
			totalOutcome += sub.StreamPrice * sub.Hours
		}
		makeRow(m, sub, t)
	}
	m.Line(1)

	m.Row(10, func() {
		m.ColSpace(7)
		m.Col(2, func() {
			m.Text("Total income:", props.Text{
				Top:   2,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Left,
			})
		})
		m.Col(3, func() {
			m.Text(strconv.FormatUint(totalIncome, 10), props.Text{
				Top:   2,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Right,
			})
		})
	})

	m.Row(10, func() {
		m.ColSpace(7)
		m.Col(2, func() {
			m.Text("Total outcome:", props.Text{
				Top:   2,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Left,
			})
		})
		m.Col(3, func() {
			m.Text(strconv.FormatUint(totalOutcome, 10), props.Text{
				Top:   2,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Right,
			})
		})
	})

	buff, err := m.Output()
	if err != nil {
		return []byte{}, err
	}
	return buff.Bytes(), nil
}

func makeRow(m pdf.Maroto, s Subscription, t string) {
	m.Row(8, func() {
		m.Col(2, func() {
			m.Text(s.StreamName, props.Text{
				Top:   2,
				Align: consts.Left,
			})
		})
		m.Col(1, func() {
			m.Text(t, props.Text{
				Top:   2,
				Align: consts.Left,
			})
		})
		m.Col(4, func() {
			m.Text(s.StartDate.Format(time.ANSIC), props.Text{
				Top:   2,
				Align: consts.Left,
			})
		})
		m.Col(4, func() {
			m.Text(s.StartDate.Format(time.ANSIC), props.Text{
				Top:   2,
				Align: consts.Left,
			})
		})
		m.Col(1, func() {
			m.Text(strconv.FormatUint(s.StreamPrice, 10), props.Text{
				Top:   2,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
	})
}
