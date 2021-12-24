package keeper_test

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tharsis/ethermint/tests"
	"github.com/tharsis/evmos/x/incentives/types"
)

func (suite *KeeperTestSuite) TestGetIncentivesGasMeters() {
	var expRes []types.GasMeter

	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"no gas meter registered",
			func() { expRes = []types.GasMeter{} },
		},
		{
			"1 gas meter registered",
			func() {
				gm := types.NewGasMeter(contract, participant, 1)
				suite.app.IncentivesKeeper.SetGasMeter(suite.ctx, gm)
				suite.Commit()

				expRes = []types.GasMeter{gm}
			},
		},
		{
			"2 gas meters registered",
			func() {
				gm := types.NewGasMeter(contract, participant, 1)
				gm2 := types.NewGasMeter(contract2, participant, 1)
				suite.app.IncentivesKeeper.SetGasMeter(suite.ctx, gm)
				suite.app.IncentivesKeeper.SetGasMeter(suite.ctx, gm2)
				suite.Commit()

				expRes = []types.GasMeter{gm, gm2}
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()
			res := suite.app.IncentivesKeeper.GetIncentivesGasMeters(suite.ctx)
			suite.Require().ElementsMatch(expRes, res, tc.name)
		})
	}
}

func (suite *KeeperTestSuite) TestGetIncentiveGasMeters() {
	var expRes []types.GasMeter

	testCases := []struct {
		name     string
		malleate func()
	}{
		// {
		// 	"no gas meter registered",
		// 	func() { expRes = []types.GasMeter{} },
		// },
		// TODO: exp result doesnt match saved as wrong hex format
		// 0x9a7de83B1691f270151A10E6c4af3754e4D09172
		// 0x7DE83b1691f270151a10E6c4Af3754e4d0917251
		{
			"1 gas meter registered",
			func() {
				gm := types.NewGasMeter(contract, participant, 1)
				suite.app.IncentivesKeeper.SetGasMeter(suite.ctx, gm)
				suite.Commit()

				expRes = []types.GasMeter{gm}
			},
			// {
			// 	"2 gas meters registered",
			// 	func() {
			// 		gm := types.NewGasMeter(contract, participant, 1)
			// 		gm2 := types.NewGasMeter(contract2, participant, 1)
			// 		suite.app.IncentivesKeeper.SetGasMeter(suite.ctx, gm)
			// 		suite.app.IncentivesKeeper.SetGasMeter(suite.ctx, gm2)
			// 		suite.Commit()

			// 		expRes = []types.GasMeter{gm, gm2}
			// 	},
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()
			res := suite.app.IncentivesKeeper.GetIncentiveGasMeters(
				suite.ctx,
				contract,
			)

			suite.Require().ElementsMatch(expRes, res, tc.name)
		})
	}
}

func (suite *KeeperTestSuite) TestGetIncentiveGasMeter() {
	expGm := types.NewGasMeter(contract, participant, 1)
	suite.app.IncentivesKeeper.SetGasMeter(suite.ctx, expGm)
	suite.Commit()

	testCases := []struct {
		name        string
		contract    common.Address
		participant common.Address
		ok          bool
	}{
		{"nil addresses", common.Address{}, common.Address{}, false},
		{"nil contract", common.Address{}, participant, false},
		{"nil paricipant", contract, common.Address{}, false},
		{"unknown contract", tests.GenerateAddress(), participant, false},
		{"unknown participant", contract, tests.GenerateAddress(), false},
		{"valid id", contract, participant, true},
	}
	for _, tc := range testCases {
		gm, found := suite.app.IncentivesKeeper.GetIncentiveGasMeter(
			suite.ctx,
			tc.contract,
			tc.participant,
		)
		if tc.ok {
			suite.Require().True(found, tc.name)
			suite.Require().Equal(expGm.CumulativeGas, gm, tc.name)
		} else {
			suite.Require().False(found, tc.name)
		}
	}
}

func (suite *KeeperTestSuite) TestDeleteGasMeter() {
	regGm := types.NewGasMeter(contract, participant, 1)
	suite.app.IncentivesKeeper.SetGasMeter(suite.ctx, regGm)
	suite.Commit()

	testCases := []struct {
		name     string
		gm       types.GasMeter
		malleate func()
		ok       bool
	}{
		{"nil incentive", types.GasMeter{}, func() {}, false},
		{"valid incentive", regGm, func() {}, true},
		{
			"detet incentive",
			regGm,
			func() {
				suite.app.IncentivesKeeper.DeleteGasMeter(suite.ctx, regGm)
			},
			false,
		},
	}
	for _, tc := range testCases {
		tc.malleate()
		gm, found := suite.app.IncentivesKeeper.GetIncentiveGasMeter(
			suite.ctx,
			common.HexToAddress(tc.gm.Contract),
			common.HexToAddress(tc.gm.Participant),
		)
		if tc.ok {
			suite.Require().True(found, tc.name)
			suite.Require().Equal(regGm.CumulativeGas, gm, tc.name)
		} else {
			suite.Require().False(found, tc.name)
		}
	}
}
