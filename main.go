package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
)

type Vote struct {
	string    `json:"creator"`
	StationId string `json:"stationId"`
	Vote      bool   `json:"vote"`
}

type VoteResult struct {
	Votes         []bool   `json:"yes"`
	Addressess    []string `json:"addressess"`
	YesCount      int      `json:"yesCount"`
	NoCount       int      `json:"noCount"`
	PercentAggree float64  `json:"percentAggree"`
	Result        bool     `json:"result"`
	Message       string   `json:"message"`
}

func main() {

	accountPath := "./accounts"
	addressPrefix := "air"

	registry, err := cosmosaccount.New(cosmosaccount.WithHome(accountPath))
	if err != nil {
		return
	}
	// create 10 accounts
	// for i := 0; i < 13; i++ {
	// 	accountName := fmt.Sprintf("account%d", i)
	// 	_, mnemonic, err := registry.Create(accountName)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	createMnemonicFile(mnemonic, accountPath, accountName) // create <accountPath>/<accountName>.mnemonic.txt file
	// }

	// create 10-10 signatures, publicKey, data, and verify all at once.
	signatureArray := make([][]byte, 13)
	publicKeyArray := make([][]byte, 13)
	dataArray := make([][]byte, 13)

	// let votes of all are these:
	votes := []bool{true, false, true, false, true, false, true, false, true, false, false, true, true}

	for i := 0; i < 13; i++ {

		accountName := fmt.Sprintf("account%d", i)
		account, err := registry.GetByName(accountName)
		addr, err := account.Address(addressPrefix)
		if err != nil {
			log.Fatal(err)
		}

		// data to be signed
		data := Vote{addr, "stationID", votes[i]}

		// Serialize data to be signed
		dataBytes, err := json.Marshal(data)
		if err != nil {
			log.Fatal("Failed to marshal data:", err)
		}

		// CREATE SIGNATURE
		signatureBytes, publicKey, err := registry.Keyring.Sign(accountName, dataBytes, 127)

		// verify signature with public key
		// verificationResult := publicKey.VerifySignature(dataBytes, signatureBytes)
		// fmt.Println("Signature verification result: ", verificationResult)

		// fmt.Println(publicKey.Type())

		signatureArray[i] = signatureBytes
		// publicKeyArray[i] = publicKey.Bytes()
		dataArray[i] = dataBytes

		interfaceRegistry := types.NewInterfaceRegistry()
		cryptocodec.RegisterInterfaces(interfaceRegistry)
		marshaler := codec.NewProtoCodec(interfaceRegistry)
		pubKeybytes, err := marshaler.MarshalInterface(publicKey)
		if err != nil {
			log.Fatalf("Failed to marshal public key: %v", err)
		}
		publicKeyArray[i] = pubKeybytes

		fmt.Println(publicKey.Address())
	}

	// verify each signature is correct
	isAllCorrectSignatures, votesArray, addressArray, err := isAllSignatureCorrect(signatureArray, publicKeyArray, dataArray)
	if !isAllCorrectSignatures {
		panic(err)
	}

	// 66% logic
	voteResult, err := voteCounts(votesArray, addressArray)
	if err != nil {
		log.Fatal(err)
	}

	data, err := json.MarshalIndent(voteResult, "", "  ") // Indent output with 2 spaces
	if err != nil {
		fmt.Println("Error marshalling:", err)
		return
	}
	fmt.Println(string(data))
	// if voteResult.Yes > 66%  voteResult.No {
	// 	yes means vrn generator is wrong,

	// 	step 1: deduct Air from the one who generated the VRN
	// 	step 2: generate new VRN

	// }else {

	// 	no means VRN generator is correct
	// 	step 1: deduct Air from the one who is the verifier of VRN
	// 	step 2: update station: vrn verified -> submit next pod : )

	// }

	// fmt.Println("Vote result: ", voteResult)
}

func createMnemonicFile(mnemonic string, accountPath string, accountName string) {
	filePath := filepath.Join(accountPath, accountName+".mnemonic.txt")
	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	_, err = f.WriteString(mnemonic)
	if err != nil {
		panic(err)
	}
	defer f.Close()
}

func isAllSignatureCorrect(signatureArray [][]byte, publicKeyArray [][]byte, dataArray [][]byte) (success bool, votesArray []bool, addressArray []string, err error) {

	pubKeysLength := len(publicKeyArray)
	if pubKeysLength != len(signatureArray) || pubKeysLength != len(dataArray) {
		errorMsg := fmt.Errorf("Length of public keys, signatures and data should be same")
		return false, nil, nil, errorMsg
	}
	if pubKeysLength == 0 {
		errorMsg := fmt.Errorf("No Votes")
		return false, nil, nil, errorMsg
	}

	// check for duplicate pubKey
	for i := 0; i < pubKeysLength; i++ {
		for j := i + 1; j < pubKeysLength; j++ {
			if string(publicKeyArray[i]) == string(publicKeyArray[j]) {
				errorMsg := fmt.Errorf("Duplicate public key found at index %d and %d", i, j)
				return false, nil, nil, errorMsg
			}
		}
	}

	// check if all address exists in the stations.tracks
	// ! TODO: get all address list from junction database
	stationsAddresses := []string{
		"33E209CAC4701B949CA5B1C07A6D427804DCEFA7",
		"11DD25A3481E5ACF4E20E895017A521566DAC4E1",
		"2F3A5D3DFE381E31CAD2DB7A1E80B260E93C4FFC",
		"2538AF3AED7D4A50110B98D5259DBB3C2FA6210B",
		"716B43234CEC050DF5C62CC27904516C47B82F25",
		"5DE13BB4A4D5E272740C9486471526C1F91FD576",
		"546ECB93D3254F852E3B410B4A3A65D089595E60",
		"D6738BCCC0A196A36EF45AB4F76F9A86506A5AB2",
		"80FFD18B65BA42636792349FE6F272523AD08580",
		"6C2445816360FBAE6F8BBBB00CB7EDCADA658C2B",
		"043D445FF1C367F5F63DC5F8C2C9777BB007B050",
		"949EFA5C7889E12487BC78CFDE7C3F25434C5F44",
		"84619427017478F15743CAE41D3AFCE02F4E96BD",
	}

	for i := 0; i < pubKeysLength; i++ {

		publicKeyBytes := publicKeyArray[i] // byte{/* your public key bytes here */}

		// Initialize a Protobuf codec for Amino
		interfaceRegistry := types.NewInterfaceRegistry()
		cryptocodec.RegisterInterfaces(interfaceRegistry)
		marshaler := codec.NewProtoCodec(interfaceRegistry)

		// Attempt to unmarshal the public key bytes into a PubKey interface
		var pubKey cryptotypes.PubKey
		err := marshaler.UnmarshalInterface(publicKeyBytes, &pubKey)
		if err != nil {
			errorMsg := fmt.Errorf("Failed to unmarshal public key: %v", err)
			return false, nil, nil, errorMsg
		}

		// check if pubKey.Address() exists in stationsAddresses
		isExists := false
		for _, address := range stationsAddresses {
			if address == pubKey.Address().String() {
				isExists = true
				break
			}
		}
		if !isExists {
			errorMsg := fmt.Errorf("Public key address %s does not exists in stationsAddresses", pubKey.Address().String())
			return false, nil, nil, errorMsg
		} else {
			fmt.Println("Public key address exists in stationsAddresses:", pubKey.Address().String())
		}

		verificationResult := pubKey.VerifySignature(dataArray[i], signatureArray[i])
		if !verificationResult {
			errorMsg := fmt.Errorf("Signature verification failed for public key address %s", pubKey.Address().String())
			return false, nil, nil, errorMsg
		}

		fmt.Println("Signature verification result: ", verificationResult)

		// unmarshal data
		var vote Vote
		err = json.Unmarshal(dataArray[i], &vote)
		if err != nil {
			errorMsg := fmt.Errorf("Failed to unmarshal data: %v", err)
			return false, nil, nil, errorMsg
		}

		votesArray = append(votesArray, vote.Vote)
		addressArray = append(addressArray, pubKey.Address().String())
	}

	return true, votesArray, addressArray, nil
}

func voteCounts(votesArray []bool, addressArray []string) (voteResult VoteResult, err error) {

	totalVotes := len(votesArray)
	yesCount := 0
	noCount := 0
	PercentAggree := 0.0
	result := false
	fmt.Println(votesArray)

	for _, vote := range votesArray {
		if vote {
			yesCount++
		} else {
			noCount++
		}
	}

	// percentagge of aggree
	PercentAggree = (float64(yesCount) / float64(totalVotes)) * 100
	// make it upto 2 decimal
	PercentAggree = float64(int(PercentAggree*10000)) / 10000

	// if yesCount > 66% of totalVotes
	if yesCount > (66*totalVotes)/100 {
		result = true
	}

	voteMsg := ""

	if result { // true
		voteMsg = "VRN generated is wrong. More than 66% votes are aggree"
	} else {
		voteMsg = "VRN generated is correct. Less than 66% votes are aggree"
	}

	voteResult = VoteResult{
		Votes:         votesArray,
		Addressess:    addressArray,
		YesCount:      yesCount,
		NoCount:       noCount,
		PercentAggree: PercentAggree,
		Result:        result,
		Message:       voteMsg,
	}

	return voteResult, nil
}
