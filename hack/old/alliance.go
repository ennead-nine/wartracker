// ScanAlliance pre-processes the given image file and scans it with tessaract
// into an alliance.Alliance struct
func (a *Alliance) ScanAlliance(img image.Image, imm scanner.ImageMaps) error {
	var d Data
	var err error

	for k, im := range imm {
		switch k {
		case "tag":
			a.Tag, err = im.ProcessImageText(img)
		case "name":
			d.Name, err = im.ProcessImageText(img)
		case "power":
			d.Power, err = im.ProcessImageInt(img)
		case "giftlevel":
			d.GiftLevel, err = im.ProcessImageInt(img)
		case "membercount":
			d.MemberCount, err = im.ProcessImageInt(img)
		default:
			return ErrInvalidMapKey
		}
		if err != nil {
			return err
		}
	}

	d.Date = time.Now().Format(time.DateOnly)

	err = a.GetByTag(a.Tag)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		fmt.Printf("A new alliance will need to be created from this data.  Please run 'wartracker-cli alliance new -o [output file]' after verifying the data\n")
	} else {
		fmt.Printf("This alliance already exists. To add the new data run 'wartracker-cli alliance add -o [output file]' to add the new data.\n")
	}

	a.DataMap[d.Date] = d

	return nil
}

// ScanAlliance pre-processes the given image file and scans it with tessaract
// into an alliance.Alliance struct
func (a *Alliance) ScanMainAlliance(img image.Image, imm scanner.ImageMaps) error {
	var d Data
	var err error

	// Setup alliance
	for k, im := range imm {
		switch k {
		case "tag":
			if cmd.Debug {
				fmt.Printf("scanning %s...\n", k)
			}
			a.Tag, err = im.ProcessImageText(img)
		case "name":
			if cmd.Debug {
				fmt.Printf("scanning %s...\n", k)
			}
			d.Name, err = im.ProcessImageText(img)
		case "power":
			if cmd.Debug {
				fmt.Printf("scanning %s...\n", k)
			}
			d.Power, err = im.ProcessImageInt(img)
		case "membercount":
			if cmd.Debug {
				fmt.Printf("scanning %s...\n", k)
			}
			tmc, err := im.ProcessImageText(img)
			if err != nil {
				return err
			}
			tmc = strings.Split(tmc, "/")[0]
			mc, err := strconv.Atoi(tmc)
			if err != nil {
				return err
			}
			d.MemberCount = mc
		default:
			return ErrInvalidMapKey
		}
		if err != nil {
			return err
		}
	}
	d.Date = time.Now().Format(time.DateOnly)

	err = a.GetByTag(a.Tag)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		fmt.Printf("A new alliance will need to be created from this data.  Please run 'wartracker-cli alliance new -o output' after verifying the data\n")
	} else {
		fmt.Printf("This alliance already exists. To add the new data run 'wartracker-cli alliance add -o output' to add the new data.\n")
		d.AllianceId = a.Id
	}

	a.DataMap[d.Date] = d

	return nil
}
