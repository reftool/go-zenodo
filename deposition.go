package gozenodo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Deposition struct {
	Created   time.Time          `json:"created"` // Creation time of deposition (in ISO8601 format)
	DOI       string             `json:"doi"`     // Digital Object Identifier (DOI)
	DOIURL    string             `json:"doi_url"` // Persistent link to your published deposition. This field is only present for published depositions
	Files     []DepositionFile   `json:"files"`   // A list of deposition files resources
	ID        int                `json:"id"`      // Deposition identifier
	Links     DepositionLinks    `json:"links"`
	Metadata  DepositionMetadata `json:"metadata"`   // A deposition metadata resource
	Modified  time.Time          `json:"modified"`   // Last modification time of deposition (in ISO8601 format).
	Owner     int                `json:"owner"`      // User identifier of the owner of the deposition.
	RecordID  int                `json:"record_id"`  // Record identifier. This field is only present for published depositions.
	RecordURL string             `json:"record_url"` // URL to public version of record for this deposition. This field is only present for published depositions.
	State     string             `json:"state"`      // inprogress, done or error
	Submitted bool               `json:"submitted"`  // True of deposition has been published, False otherwise.
	Title     string             `json:"title"`      // Title of deposition (automatically set from metadata). Defaults to empty string.
}

type DepositionLinks struct {
	Bucket          string `json:"bucket"`
	Discard         string `json:"discard"`
	Edit            string `json:"edit"`
	Files           string `json:"files"`
	HTML            string `json:"html"`
	LatestDraft     string `json:"latest_draft"`
	LatestDraftHTML string `json:"latest_draft_html"`
	Publish         string `json:"publish"`
	Self            string `json:"self"`
}

type DepositionMetadata struct {
	/*
		Controlled vocabulary:
		* publication: Publication
		* poster: Poster
		* presentation: Presentation
		* dataset: Dataset
		* image: Image
		* video: Video/Audio
		* software: Software
		* lesson: Lesson
		* physicalobject: Physical object
		* other: Other
	*/
	UploadType string `json:"upload_type"`

	/*
		Controlled vocabulary:
		* annotationcollection: Annotation collection
		* book: Book
		* section: Book section
		* conferencepaper: Conference paper
		* datamanagementplan: Data management plan
		* article: Journal article
		* patent: Patent
		* preprint: Preprint
		* deliverable: Project deliverable
		* milestone: Project milestone
		* proposal: Proposal
		* report: Report
		* softwaredocumentation: Software documentation
		* taxonomictreatment: Taxonomic treatment
		* technicalnote: Technical note
		* thesis: Thesis
		* workingpaper: Working paper
		* other: Other
	*/
	PublicationType string `json:"publication_type"`

	/*
		Controlled vocabulary:
		* figure: Figure
		* plot: Plot
		* drawing: Drawing
		* diagram: Diagram
		* photo: Photo
		* other: Other
	*/
	ImageType string `json:"image_type"`

	PublicationDate string `json:"publication_date"` // Date of publication in ISO8601 format (YYYY-MM-DD). Defaults to current date.
	Title           string `json:"title"`            // Title of deposition.

	Creators    []DepositionCreator `json:"creators"`     // The creators/authors of the deposition.
	Description string              `json:"description"`  // Abstract or description for deposition. (allows HTML)
	AccessRight string              `json:"access_right"` // open, embargoed, restricted, closed

	/*
		The selected license applies to all files in this deposition, but not to the metadata which is licensed under
		Creative Commons Zero. You can find the available license IDs via our /api/licenses endpoint. Defaults to cc-zero
		for datasets and cc-by for everything else.
	*/
	License string `json:"license"`

	EmbargoDate string `json:"embargo_date"` // When the deposited files will be made automatically made publicly available by the system. Defaults to current date.

	/*
		Specify the conditions under which you grant users access to the files in your upload. User requesting access
		will be asked to justify how they fulfil the conditions. Based on the justification, you decide who to grant/deny access.
		You are not allowed to charge users for granting access to data hosted on Zenodo.
	*/
	AccessConditions string `json:"access_conditions"`

	/*
		Digital Object Identifier. Did a publisher already assign a DOI to your deposited files? If not, leave the field
		empty and we will register a new DOI for you when you publish. A DOI allow others to easily and unambiguously cite
		your deposition.
	*/
	DOI string `json:"doi"`

	/*
		Set to true, to reserve a Digital Object Identifier (DOI). The DOI is automatically generated by our system and
		cannot be changed. Also, The DOI is not registered with DataCite until you publish your deposition, and thus
		cannot be used before then. Reserving a DOI is useful, if you need to include it in the files you upload, or if
		you need to provide a dataset DOI to your publisher but not yet publish your dataset. The response from the REST
		API will include the reserved DOI.
	*/
	PreserveDOI bool `json:"preserve_doi"`

	Keywords []string `json:"keywords"` // Free form keywords for this deposition.

	Notes string `json:"notes"` // Additional notes.

	/*
		Persistent identifiers of related publications and datasets. Supported identifiers include: DOI, Handle, ARK,
		PURL, ISSN, ISBN, PubMed ID, PubMed Central ID, ADS Bibliographic Code, arXiv, Life Science Identifiers (LSID),
		EAN-13, ISTC, URNs and URLs.
	*/
	RelatedIdentifiers []DepositionIdentifier `json:"related_identifiers"`

	Contributors []DepositionContributor `json:"contributors"` // The contributors of the deposition (e.g. editors, data curators, etc.).
	References   []string                `json:"references"`   // ListDepositions of references

	/*
		ListDepositions of communities you wish the deposition to appear. The owner of the community will be notified, and can
		either accept or reject your request.
	*/
	Communities []DepositionCommunity `json:"communities"`

	Grants []DepositionGrant `json:"grants"` // ListDepositions of OpenAIRE-supported grants, which have funded the research for this deposition.

	JournalTitle  string `json:"journal_title"`  // Journal title, if deposition is a published article.
	JournalVolume string `json:"journal_volume"` // Journal volume, if deposition is a published article.
	JournalIssue  string `json:"journal_issue"`  // Journal issue, if deposition is a published article.
	JournalPages  string `json:"journal_pages"`  // Journal pages, if deposition is a published article.

	ConferenceTitle       string `json:"conference_title"`        // Title of conference (e.g. 20th International Conference on Computing in High Energy and Nuclear Physics).
	ConferenceAcronym     string `json:"conference_acronym"`      // Acronym of conference (e.g. CHEP'13).
	ConferenceDates       string `json:"conference_dates"`        // Dates of conference (e.g. 14-18 October 2013). Conference title or acronym must also be specified if this field is specified.
	ConferencePlace       string `json:"conference_place"`        // Place of conference in the format city, country (e.g. Amsterdam, The Netherlands). Conference title or acronym must also be specified if this field is specified.
	ConferenceURL         string `json:"conference_url"`          // URL of conference (e.g. http://www.chep2013.org/).
	ConferenceSession     string `json:"conference_session"`      // Number of session within the conference (e.g. VI).
	ConferenceSessionPart string `json:"conference_session_part"` // Number of part within a session (e.g. 1).

	ImprintPublisher string `json:"imprint_publisher"` // Publisher of a book/report/chapter
	ImprintISBN      string `json:"imprint_isbn"`      // ISBN of a book/report
	ImprintPlace     string `json:"imprint_place"`     // Place of publication of a book/report/chapter in the format city, country.

	PartOfTitle string `json:"partof_title"` // Title of book for chapters
	PartOfPages string `json:"partof_pages"` // Pages numbers of book

	ThesisSupervisors []DepositionThesisSupervisor `json:"thesis_supervisors"` // Supervisors of the thesis.
	ThesisUniversity  string                       `json:"thesis_university"`  // Awarding university of thesis.

	Subjects []DepositionSubject `json:"subjects"` // Specify subjects from a taxonomy or controlled vocabulary. Each term must be uniquely identified (e.g. a URL). For free form text, use the keywords field

	Version   string               `json:"version"`   // Version of the resource. Any string will be accepted, however the suggested format is a semantically versioned tag (see more details on semantic versioning at semver.org)
	Language  string               `json:"language"`  // Specify the main language of the record as ISO 639-2 or 639-3 code, see Library of Congress ISO 639 codes list.
	Locations []DepositionLocation `json:"locations"` // ListDepositions of locations
	Dates     []DepositionDate     `json:"dates"`     // ListDepositions of date intervals
	Method    string               `json:"method"`    // The methodology employed for the study or research.
}

type DepositionCreator struct {
	Name        string `json:"name"`
	Affiliation string `json:"affiliation"`
	ORCID       string `json:"orcid"`
	GND         string `json:"gnd"`
}

type DepositionIdentifier struct {
	Identifier   string `json:"identifier"`
	Relation     string `json:"relation"`
	ResourceType string `json:"resource_type"`
}

type DepositionContributor struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Affiliation string `json:"affiliation"`
	ORCID       string `json:"orcid"`
	GND         string `json:"gnd"`
}

type DepositionCommunity struct {
	Identifier string `json:"identifier"`
}

type DepositionGrant struct {
	ID string `json:"id"`
}

type DepositionThesisSupervisor struct {
	Name        string `json:"name"`
	Affiliation string `json:"affiliation"`
	ORCID       string `json:"orcid"`
	GND         string `json:"gnd"`
}

type DepositionSubject struct {
	Term       string `json:"term"`
	Identifier string `json:"identifier"`
	Scheme     string `json:"scheme"`
}

type DepositionLocation struct {
	Lat         float64 `json:"lat"`
	Long        float64 `json:"long"`
	Place       string  `json:"place"`
	Description string  `json:"description"`
}

type DepositionDate struct {
	Start       string `json:"start"`
	End         string `json:"end"`
	Type        string `json:"type"` // Collected, Valid, Withdrawn
	Description string `json:"description"`
}

type DepositionFile struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Filesize int    `json:"filesize"` // size of file in bytes
	Checksum string `json:"checksum"` // MD5 checksum of file, computed by zenodo. This allows you to check the integrity of the uploaded file.
}

type DepositionFileUpload struct {
	Created      time.Time                 `json:"created"`
	Updated      time.Time                 `json:"updated"`
	VersionID    string                    `json:"version_id"`
	Key          string                    `json:"key"`
	Size         int                       `json:"size"`
	Mimetype     string                    `json:"mimetype"`
	Checksum     string                    `json:"checksum"`
	IsHead       bool                      `json:"is_head"`
	DeleteMarker bool                      `json:"delete_marker"`
	Links        DepositionFileUploadLinks `json:"links"`
}

type DepositionFileUploadLinks struct {
	Self    string `json:"self"`
	Version string `json:"version"`
	Uploads string `json:"uploads"`
}

func CreateDeposition() (*Deposition, error) {
	if Token == "" {
		return nil, errors.New("token not set. use SetAccessToken")
	}

	var response *Deposition

	zenodoBaseURL := SandboxURL
	if !SandboxMode {
		zenodoBaseURL = ProdURL
	}

	requestURL := fmt.Sprintf("%s/api/deposit/depositions", zenodoBaseURL)

	req, err := http.NewRequest("POST", requestURL, strings.NewReader("{}"))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("couldn't decode response from zenodo")
	}

	if res.StatusCode != http.StatusCreated {
		return nil, errors.New("couldn't create deposition (" + string(data) + ")")
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func GetDeposition(id int) (*Deposition, error) {
	if Token == "" {
		return nil, errors.New("token not set. use SetAccessToken")
	}

	var response *Deposition

	zenodoBaseURL := SandboxURL
	if !SandboxMode {
		zenodoBaseURL = ProdURL
	}

	requestURL := fmt.Sprintf("%s/api/deposit/depositions/%d", zenodoBaseURL, id)

	req, err := http.NewRequest("GET", requestURL, strings.NewReader("{}"))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("couldn't decode response from zenodo")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("couldn't get deposition (" + string(data) + ")")
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func DeleteDeposition(id int) error {
	if Token == "" {
		return errors.New("token not set. use SetAccessToken")
	}

	zenodoBaseURL := SandboxURL
	if !SandboxMode {
		zenodoBaseURL = ProdURL
	}

	requestURL := fmt.Sprintf("%s/api/deposit/depositions/%d", zenodoBaseURL, id)

	req, err := http.NewRequest("DELETE", requestURL, strings.NewReader("{}"))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusNoContent {
		fmt.Println(res.StatusCode)
		return errors.New("couldn't delete deposition")
	}

	return nil
}

func ListDepositions() ([]*Deposition, error) {
	if Token == "" {
		return nil, errors.New("token not set. use SetAccessToken")
	}

	var response []*Deposition

	zenodoBaseURL := SandboxURL
	if !SandboxMode {
		zenodoBaseURL = ProdURL
	}

	requestURL := fmt.Sprintf("%s/api/deposit/depositions", zenodoBaseURL)

	req, err := http.NewRequest("GET", requestURL, strings.NewReader("{}"))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("couldn't decode response from zenodo")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("couldn't get depositions")
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response, nil
}

type UpdateDepositionRequest struct {
	Metadata *Deposition `json:"metadata"`
}

func UpdateDeposition(d *Deposition) (*Deposition, error) {
	if Token == "" {
		return nil, errors.New("token not set. use SetAccessToken")
	}

	var response *Deposition

	zenodoBaseURL := SandboxURL
	if !SandboxMode {
		zenodoBaseURL = ProdURL
	}

	// convert deposition to json string
	newRequestData := UpdateDepositionRequest{
		Metadata: d,
	}

	jsonData, err := json.Marshal(newRequestData)
	if err != nil {
		return nil, err
	}

	requestURL := fmt.Sprintf("%s/api/deposit/depositions/%d", zenodoBaseURL, d.ID)
	req, err := http.NewRequest("PUT", requestURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("couldn't decode response from zenodo")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("couldn't updatedeposition (" + string(data) + ")")
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func UploadFile(bucketID, fileName, path string) (*DepositionFileUpload, error) {
	if Token == "" {
		return nil, errors.New("token not set. use SetAccessToken")
	}

	var response *DepositionFileUpload

	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	requestURL := fmt.Sprintf("%s/%s", bucketID, fileName)
	req, err := http.NewRequest("PUT", requestURL, strings.NewReader(string(fileData)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("couldn't decode response from zenodo")
	}

	if res.StatusCode != http.StatusCreated {
		return nil, errors.New("couldn't upload file")
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response, nil
}
