package pontuz

type NotarListing struct {
	EntityType        string `json:"entityType"`
	EntityID          string `json:"entityId"`
	ID                string `json:"id"`
	OriginalEstateID  string `json:"originalEstateId"`
	Type              string `json:"type"`
	IndexType         string `json:"indexType"`
	AgentID           string `json:"agentId"`
	AdditionalAgentID string `json:"additionalAgentId"`
	CustomerID        string `json:"customerId"`
	OfficeID          string `json:"officeId"`
	ObjectSettings    struct {
		ShowPrice              bool `json:"showPrice"`
		ShowAddress            bool `json:"showAddress"`
		ShowAsComing           bool `json:"showAsComing"`
		ShowAsPreview          bool `json:"showAsPreview"`
		ShowAsReferenceHousing bool `json:"showAsReferenceHousing"`
		ShowAsSoonForSale      bool `json:"showAsSoonForSale"`
		ShowAsTodaysHousing    bool `json:"showAsTodaysHousing"`
	} `json:"objectSettings"`
	UpdatedAt         string `json:"updatedAt"`
	UpdatedRemotelyAt string `json:"updatedRemotelyAt"`
	Price             struct {
		Currency                        string `json:"currency"`
		Type                            string `json:"type"`
		StartingPrice                   int    `json:"startingPrice"`
		StartingLivingSpacePricePerArea int    `json:"startingLivingSpacePricePerArea"`
	} `json:"price,omitempty"`
	Address          string  `json:"address"`
	County           string  `json:"county"`
	Area             string  `json:"area"`
	Municipality     string  `json:"municipality"`
	ZipCode          string  `json:"zipCode"`
	City             string  `json:"city"`
	Longitude        float64 `json:"longitude"`
	Latitude         float64 `json:"latitude"`
	CountryCode      string  `json:"countryCode"`
	LivingSpaceArea  int     `json:"livingSpaceArea"`
	NumberOfRooms    int     `json:"numberOfRooms"`
	AssignmentStatus string  `json:"assignmentStatus"`
	Active           bool    `json:"active"`
	Inspected        bool    `json:"inspected"`
	InspectedAlt     bool    `json:"inspected_alt"`
	Images           []struct {
		File     string `json:"file"`
		Order    int    `json:"order"`
		Category string `json:"category"`
	} `json:"images"`
	ShowOneViewingTime bool `json:"showOneViewingTime"`
	Viewings           []struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"viewings"`
	Bids       []any  `json:"bids"`
	BidSetting string `json:"bidSetting"`
	Interior   struct {
		NumberOfRooms int    `json:"numberOfRooms"`
		Description   string `json:"description"`
		Rooms         []struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		} `json:"rooms"`
		CompiledDescription bool `json:"compiledDescription"`
	} `json:"interior,omitempty"`
	Other struct {
		AreaSource        string `json:"areaSource"`
		AreaSourceComment string `json:"areaSourceComment"`
		OtherSpaceArea    int    `json:"otherSpaceArea"`
		LocalIncomeTax    string `json:"localIncomeTax"`
		Broadband         string `json:"broadband"`
	} `json:"other,omitempty"`
	Residence struct {
		AreaID          string `json:"areaId"`
		Area            string `json:"area"`
		Municipality    string `json:"municipality"`
		LivingSpaceArea struct {
			LivingSpaceArea   int    `json:"livingSpaceArea"`
			AreaSource        string `json:"areaSource"`
			AreaSourceComment string `json:"areaSourceComment"`
		} `json:"livingSpaceArea"`
		OtherSpaceArea int `json:"otherSpaceArea"`
		NumberOfRooms  int `json:"numberOfRooms"`
		TvAndBroadband struct {
			Broadband string `json:"broadband"`
		} `json:"tvAndBroadband"`
		WaterAndDrain string `json:"waterAndDrain"`
		Tax           struct {
			TypeCode                 string `json:"typeCode"`
			TaxAssessmentYear        int    `json:"taxAssessmentYear"`
			TaxFee                   int    `json:"taxFee"`
			LandValue                int    `json:"landValue"`
			TotalAssessedValue       int    `json:"totalAssessedValue"`
			PreliminaryAssessedValue bool   `json:"preliminaryAssessedValue"`
		} `json:"tax"`
		Mortgage struct {
		} `json:"mortgage"`
		Enrollment struct {
			PreferentialAndCommunity string `json:"preferentialAndCommunity"`
		} `json:"enrollment"`
		EnergyDeclaration struct {
			CompletedDescription string `json:"completedDescription"`
		} `json:"energyDeclaration"`
		DisposalForm            string `json:"disposalForm"`
		PropertyType            string `json:"propertyType"`
		BuildingType            string `json:"buildingType"`
		PropertyUnitDesignation string `json:"propertyUnitDesignation"`
		Tenure                  string `json:"tenure"`
		PlotArea                int    `json:"plotArea"`
		Patio                   string `json:"patio"`
		BuildingYear            string `json:"buildingYear"`
		Construction            struct {
			Facade     string `json:"facade"`
			Foundation string `json:"foundation"`
			Frame      string `json:"frame"`
			Windows    string `json:"windows"`
			Roof       string `json:"roof"`
		} `json:"construction"`
		Heating     string `json:"heating"`
		PlotParking string `json:"plotParking"`
		Leasehold   struct {
		} `json:"leasehold"`
		Lease struct {
		} `json:"lease"`
	} `json:"residence,omitempty"`
	NotarPlus    bool   `json:"notarPlus"`
	EstateNumber string `json:"estateNumber"`
	Plot         struct {
		Area               int    `json:"area"`
		BuildingPermission any    `json:"buildingPermission"`
		Other              string `json:"other"`
		Parking            string `json:"parking"`
		Patio              string `json:"patio"`
		PlotType           string `json:"plotType"`
	} `json:"plot"`
	Other0 struct {
		AreaSource     string `json:"areaSource"`
		OtherSpaceArea int    `json:"otherSpaceArea"`
		LocalIncomeTax string `json:"localIncomeTax"`
	} `json:"other,omitempty"`
	Residence0 struct {
		AreaID          string `json:"areaId"`
		Area            string `json:"area"`
		Municipality    string `json:"municipality"`
		LivingSpaceArea struct {
			LivingSpaceArea int    `json:"livingSpaceArea"`
			AreaSource      string `json:"areaSource"`
		} `json:"livingSpaceArea"`
		OtherSpaceArea int `json:"otherSpaceArea"`
		NumberOfRooms  int `json:"numberOfRooms"`
		TvAndBroadband struct {
		} `json:"tvAndBroadband"`
		WaterAndDrain string `json:"waterAndDrain"`
		Tax           struct {
			TypeCode                 string `json:"typeCode"`
			TaxAssessmentYear        int    `json:"taxAssessmentYear"`
			TaxFee                   int    `json:"taxFee"`
			LandValue                int    `json:"landValue"`
			BuildingValue            int    `json:"buildingValue"`
			TotalAssessedValue       int    `json:"totalAssessedValue"`
			PreliminaryAssessedValue bool   `json:"preliminaryAssessedValue"`
			ValueYear                int    `json:"valueYear"`
		} `json:"tax"`
		Mortgage struct {
			NumberOfMortgages   int `json:"numberOfMortgages"`
			TotalMortgageAmount int `json:"totalMortgageAmount"`
		} `json:"mortgage"`
		Enrollment struct {
			PlanRegulations          string `json:"planRegulations"`
			PreferentialAndCommunity string `json:"preferentialAndCommunity"`
		} `json:"enrollment"`
		EnergyDeclaration struct {
			CompletedDescription string `json:"completedDescription"`
			Class                string `json:"class"`
			Consumption          string `json:"consumption"`
			DeclarationDate      string `json:"declarationDate"`
			PrimaryEnergyNumber  string `json:"primaryEnergyNumber"`
		} `json:"energyDeclaration"`
		DisposalForm            string `json:"disposalForm"`
		PropertyType            string `json:"propertyType"`
		BuildingType            string `json:"buildingType"`
		PropertyUnitDesignation string `json:"propertyUnitDesignation"`
		Tenure                  string `json:"tenure"`
		PlotArea                int    `json:"plotArea"`
		BuildingYear            string `json:"buildingYear"`
		Construction            struct {
			Facade     string `json:"facade"`
			Foundation string `json:"foundation"`
			Frame      string `json:"frame"`
			Windows    string `json:"windows"`
			Roof       string `json:"roof"`
		} `json:"construction"`
		Renovations string `json:"renovations"`
		Heating     string `json:"heating"`
		Leasehold   struct {
		} `json:"leasehold"`
		Lease struct {
		} `json:"lease"`
	} `json:"residence,omitempty"`
	Interior0 struct {
		NumberOfRooms       int   `json:"numberOfRooms"`
		Rooms               []any `json:"rooms"`
		CompiledDescription bool  `json:"compiledDescription"`
	} `json:"interior,omitempty"`
	Other1 struct {
		AreaSource     string `json:"areaSource"`
		LocalIncomeTax string `json:"localIncomeTax"`
	} `json:"other,omitempty"`
	Residence1 struct {
		AreaID          string `json:"areaId"`
		Area            string `json:"area"`
		Municipality    string `json:"municipality"`
		LivingSpaceArea struct {
			LivingSpaceArea int    `json:"livingSpaceArea"`
			AreaSource      string `json:"areaSource"`
		} `json:"livingSpaceArea"`
		NumberOfRooms  int `json:"numberOfRooms"`
		TvAndBroadband struct {
		} `json:"tvAndBroadband"`
		WaterAndDrain string `json:"waterAndDrain"`
		Tax           struct {
			TypeCode                 string `json:"typeCode"`
			TaxAssessmentYear        int    `json:"taxAssessmentYear"`
			TaxFee                   int    `json:"taxFee"`
			LandValue                int    `json:"landValue"`
			TotalAssessedValue       int    `json:"totalAssessedValue"`
			PreliminaryAssessedValue bool   `json:"preliminaryAssessedValue"`
		} `json:"tax"`
		Mortgage struct {
			NumberOfMortgages   int `json:"numberOfMortgages"`
			TotalMortgageAmount int `json:"totalMortgageAmount"`
		} `json:"mortgage"`
		Enrollment struct {
			PlanRegulations string `json:"planRegulations"`
		} `json:"enrollment"`
		EnergyDeclaration struct {
			CompletedDescription string `json:"completedDescription"`
			Class                string `json:"class"`
			DeclarationDate      string `json:"declarationDate"`
			PrimaryEnergyNumber  string `json:"primaryEnergyNumber"`
		} `json:"energyDeclaration"`
		DisposalForm            string `json:"disposalForm"`
		PropertyType            string `json:"propertyType"`
		BuildingType            string `json:"buildingType"`
		PropertyUnitDesignation string `json:"propertyUnitDesignation"`
		Tenure                  string `json:"tenure"`
		PlotArea                int    `json:"plotArea"`
		BuildingYear            string `json:"buildingYear"`
		Construction            struct {
			Facade      string `json:"facade"`
			Foundation  string `json:"foundation"`
			Frame       string `json:"frame"`
			Windows     string `json:"windows"`
			Roof        string `json:"roof"`
			Ventilation string `json:"ventilation"`
		} `json:"construction"`
		Heating   string `json:"heating"`
		Leasehold struct {
		} `json:"leasehold"`
		Lease struct {
		} `json:"lease"`
	} `json:"residence,omitempty"`
	Price0 struct {
		Currency                        string `json:"currency"`
		StartingPrice                   int    `json:"startingPrice"`
		StartingLivingSpacePricePerArea int    `json:"startingLivingSpacePricePerArea"`
	} `json:"price,omitempty"`
	Other2 struct {
		AreaSource     string `json:"areaSource"`
		OtherSpaceArea int    `json:"otherSpaceArea"`
		LocalIncomeTax string `json:"localIncomeTax"`
		Broadband      string `json:"broadband"`
	} `json:"other,omitempty"`
	Residence2 struct {
		AreaID          string `json:"areaId"`
		Area            string `json:"area"`
		Municipality    string `json:"municipality"`
		LivingSpaceArea struct {
			LivingSpaceArea int    `json:"livingSpaceArea"`
			AreaSource      string `json:"areaSource"`
		} `json:"livingSpaceArea"`
		OtherSpaceArea int `json:"otherSpaceArea"`
		NumberOfRooms  int `json:"numberOfRooms"`
		TvAndBroadband struct {
			Broadband string `json:"broadband"`
		} `json:"tvAndBroadband"`
		WaterAndDrain string `json:"waterAndDrain"`
		Tax           struct {
			TypeCode                 string `json:"typeCode"`
			TaxAssessmentYear        int    `json:"taxAssessmentYear"`
			TaxFee                   int    `json:"taxFee"`
			LandValue                int    `json:"landValue"`
			BuildingValue            int    `json:"buildingValue"`
			TotalAssessedValue       int    `json:"totalAssessedValue"`
			PreliminaryAssessedValue bool   `json:"preliminaryAssessedValue"`
			ValueYear                int    `json:"valueYear"`
		} `json:"tax"`
		Mortgage struct {
			NumberOfMortgages   int `json:"numberOfMortgages"`
			TotalMortgageAmount int `json:"totalMortgageAmount"`
		} `json:"mortgage"`
		Enrollment struct {
			PlanRegulations          string `json:"planRegulations"`
			PreferentialAndCommunity string `json:"preferentialAndCommunity"`
		} `json:"enrollment"`
		EnergyDeclaration struct {
			CompletedDescription string `json:"completedDescription"`
			Class                string `json:"class"`
			Consumption          string `json:"consumption"`
			DeclarationDate      string `json:"declarationDate"`
			PrimaryEnergyNumber  string `json:"primaryEnergyNumber"`
		} `json:"energyDeclaration"`
		DisposalForm            string `json:"disposalForm"`
		PropertyType            string `json:"propertyType"`
		BuildingType            string `json:"buildingType"`
		PropertyUnitDesignation string `json:"propertyUnitDesignation"`
		Tenure                  string `json:"tenure"`
		PlotArea                int    `json:"plotArea"`
		PlotOtherDescription    string `json:"plotOtherDescription"`
		Patio                   string `json:"patio"`
		BuildingYear            string `json:"buildingYear"`
		Construction            struct {
			Facade      string `json:"facade"`
			Foundation  string `json:"foundation"`
			Frame       string `json:"frame"`
			Windows     string `json:"windows"`
			Beam        string `json:"beam"`
			Roof        string `json:"roof"`
			Ventilation string `json:"ventilation"`
		} `json:"construction"`
		Renovations string `json:"renovations"`
		Heating     string `json:"heating"`
		PlotParking string `json:"plotParking"`
		Leasehold   struct {
		} `json:"leasehold"`
		Lease struct {
		} `json:"lease"`
	} `json:"residence,omitempty"`
}
