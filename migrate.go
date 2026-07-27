package migrate

// ... existing imports ...

func (m *Migrate) apply(migration *Migration) error {
	// ... existing code ...
	if m.config.Transaction {
		if err := m.db.Begin(); err != nil {
			return err
		}
		defer m.db.Rollback()
	}

	if err := migration.Up(m.db); err != nil {
		// If transaction is supported, rollback will happen via defer
		// Only mark dirty if we are not in a transaction or if rollback fails
		if !m.config.Transaction {
			m.setDirty(true)
		}
		return err
	}

	if m.config.Transaction {
		if err := m.db.Commit(); err != nil {
			m.setDirty(true)
			return err
		}
	}

	m.setVersion(migration.Version, false)
	return nil
}