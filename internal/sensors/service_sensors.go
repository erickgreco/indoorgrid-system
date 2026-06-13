package sensors

import (
	"context"

	"github.com/erickgreco/indoorgrid-system/pkg/logger"
	"github.com/erickgreco/indoorgrid-system/pkg/syserrors"
	"github.com/google/uuid"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) SaveBME680(ctx context.Context, reading BME680Reading) error {

	if reading.TemperatureCelcius < minTemperature || reading.TemperatureCelcius > maxTemperature {
		return logger.Error(logger.ReadingErr, syserrors.ErrInvalidTemp, logger.Reading, reading.TemperatureCelcius)
	}

	if reading.HumidityPercent < minHumidity || reading.HumidityPercent > maxHumidity {
		return logger.Error(logger.ReadingErr, syserrors.ErrInvalidHum, logger.Reading, reading.HumidityPercent)
	}

	if reading.PressureHpa < minPressureHpa || reading.PressureHpa > maxPressureHpa {
		return logger.Error(logger.ReadingErr, syserrors.ErrInvalidPress, logger.Reading, reading.PressureHpa)
	}

	if reading.GasResistanceOhms < minGasResistanceOhms || reading.GasResistanceOhms > maxGasResistanceOhms {
		return logger.Error(logger.ReadingErr, syserrors.ErrInvalidGas, logger.Reading, reading.GasResistanceOhms)
	}

	reading.ID = uuid.New()

	cycleID, err := s.repo.getActiveCycleID(ctx)
	if err != nil {
		return logger.Error(logger.ActiveCycleErr, err)
	}
	if cycleID != nil {
		reading.GrowCycleID = cycleID
	}

	return s.repo.insertBME680(ctx, reading)
}

func (s *Service) SaveBH1750(ctx context.Context, reading BH1750Reading) error {
	if reading.IlluminanceLux < minIlluminanceLux || reading.IlluminanceLux > maxIlluminanceLux {
		return logger.Error(logger.ReadingErr, syserrors.ErrInvalidLux, logger.Reading, reading.IlluminanceLux)
	}

	reading.ID = uuid.New()

	cycleID, err := s.repo.getActiveCycleID(ctx)
	if err != nil {
		return logger.Error(logger.ActiveCycleErr, err)
	}
	if cycleID != nil {
		reading.GrowCycleID = cycleID
	}

	return s.repo.insertBH1750(ctx, reading)
}

func (s *Service) SaveSoil(ctx context.Context, reading SoilReading) error {
	if reading.MoisturePercent < minSoilPercent || reading.MoisturePercent > maxSoilPercent {
		return logger.Error(logger.ReadingErr, syserrors.ErrInvalidSoil, logger.Reading, reading.MoisturePercent)
	}

	reading.ID = uuid.New()

	cycleID, err := s.repo.getActiveCycleID(ctx)
	if err != nil {
		return logger.Error(logger.ActiveCycleErr, err)
	}
	if cycleID != nil {
		reading.GrowCycleID = cycleID
	}

	return s.repo.insertSoil(ctx, reading)
}
