package controllers

import (
	"backend/src/db"
	"backend/src/models"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"sort"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type UserService struct {
	DbAdapter *db.DbAdapter
}

func NewUserService(dbAdapter *db.DbAdapter) *UserService {
	return &UserService{DbAdapter: dbAdapter}
}

func (u UserService) CreateUser(ctx context.Context, r *http.Request) (models.Response, error) {
	file, _, err := r.FormFile("paymentImg")

	if err != nil {
		log.Println(err)
		return models.Response{Message: "Error uploading file", Success: false}, err
	}

	defer file.Close()

	url, success := u.FileUpload(context.Background(), file)

	if !success {
		return models.Response{Message: "Error uploading file", Success: false}, nil
	}

	userInput, err := u.GetUserInfo(r)
	userInput.PaymentImg = url

	log.Println(userInput)
	if err != nil {
		return models.Response{Message: "Error getting user info", Success: false}, err
	}

	isValid, message := u.ValidateUserInput(*userInput)

	if !isValid {
		return models.Response{Message: message, Success: false}, nil
	}
	userID, err := u.DbAdapter.CreateUser(ctx, *userInput)

	if err != nil {
		return models.Response{Message: "Error creating user", Success: false, Error: err.Error()}, err
	}

	go func() {
		isSent := u.SendEmail(*userInput)
		if !isSent {
			fmt.Println("Failed to send email:", err)
		} else {
			fmt.Println("Email sent successfully to")
		}
	}()
	return models.Response{Message: "User created successfully", Data: userID, Success: true}, nil

}

func (u UserService) SendEmail(user models.UserInput) bool {
	from := os.Getenv("BACKEND_MAIL_USER")
	password := os.Getenv("BACKEND_MAIL_PASSWORD")
	host := os.Getenv("BACKEND_MAIL_HOST")
	to := user.Email

	log.Println(from, password, host, to)

	auth := smtp.PlainAuth("", from, password, host)

	emailTemplate := u.GetEmail(user.Name)

	msg := []byte("From: " + from + "\r\n" + "To: " + to + "\r\n" +
		"Subject: Welcome to LinuxDiary 5.0\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
		emailTemplate)

	err := smtp.SendMail(host+":587", auth, from, []string{to}, msg)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func (u UserService) GetUserInfo(r *http.Request) (*models.UserInput, error) {
	err := r.ParseMultipartForm(10 << 20)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	userInput := models.UserInput{
		Name:          r.FormValue("name"),
		Email:         r.FormValue("email"),
		Phone:         r.FormValue("phone"),
		TransactionId: r.FormValue("transactionId"),
		CollegeName:   r.FormValue("collegeName"),
		YearOfStudy:   r.FormValue("yearOfStudy"),
		Branch:        r.FormValue("branch"),
		IsDualBooted:  r.FormValue("isDualBooted") == "true",
		ReferralCode:  r.FormValue("referralCode"),
	}

	return &userInput, nil
}

func (u UserService) ValidateUserInput(userInput models.UserInput) (bool, string) {
	if userInput.Name == "" {
		return false, "Name is required"
	}
	if userInput.Email == "" {
		return false, "Email is required"
	}
	if userInput.Phone == "" {
		return false, "Phone is required"
	}
	if userInput.TransactionId == "" {
		return false, "TransactionId is required"
	}
	if userInput.CollegeName == "" {
		return false, "CollegeName is required"
	}
	if userInput.YearOfStudy == "" {
		return false, "YearOfStudy is required"
	}
	if userInput.Branch == "" {
		return false, "Branch is required"
	}
	if userInput.PaymentImg == "" {
		return false, "PaymentImg is required"
	}
	return true, ""
}

func (u UserService) GetEmail(name string) string {
	return `<!DOCTYPE html>
    <html lang="en">
    <head>
      <meta charset="UTF-8" />
      <meta http-equiv="X-UA-Compatible" content="IE=edge" />
      <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      <link rel="preconnect" href="https://fonts.googleapis.com" />
      <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
      <link
        href="https://fonts.googleapis.com/css2?family=Poppins:ital,wght@0,400;0,500;0,600;1,400;1,500&display=swap"
        rel="stylesheet"
      />
      <title>LinuxDiary 5.0</title>
    </head>
    
    <body style="font-family: 'Poppins', sans-serif">
      <div>
        <u></u>
    
        <div
          style="
            text-align: center;
            margin: 0;
            padding-top: 10px;
            padding-bottom: 10px;
            padding-left: 0;
            padding-right: 0;
            background-color: #f2f4f6;
            color: #000000;
          "
          align="center"
        >
          <div style="text-align: center">
            <table
              align="center"
              style="
                text-align: center;
                vertical-align: middle;
                width: 600px;
                max-width: 600px;
              "
              width="600"
            >
              <tbody>
                <tr>
                  <td
                    style="width: 596px; vertical-align: middle"
                    width="596"
                  ></td>
                </tr>
              </tbody>
            </table>
    
            <img
              style="
                width: 600px;
                max-width: 600px;
                height: 350px;
                max-height: 350px;
                text-align: center;
              "
              alt="LinuxDiary 5.0 image"
              src="https://res.cloudinary.com/dbqq0kjga/image/upload/v1691067553/LINUX-Diary_j7sq1s.png"
              align="center"
              width="600"
              height="350"
              class="CToWUd a6T"
              data-bit="iit"
              tabindex="0"
            />
            <div
              class="a6S"
              dir="ltr"
              style="opacity: 0.01; left: 552px; top: 501.5px"
            >
              <div
                id=":155"
                class="T-I J-J5-Ji aQv T-I-ax7 L3 a5q"
                role="button"
                tabindex="0"
                aria-label="Download attachment "
                jslog="91252; u014N:cOuCgd,Kr2w4b,xr6bB; 4:WyIjbXNnLWY6MTc2MjU0MTQxMTA0MjYyMTM2NyIsbnVsbCxbXV0."
                data-tooltip-class="a1V"
                data-tooltip="Download"
              >
                <div class="akn">
                  <div class="aSK J-J5-Ji aYr"></div>
                </div>
              </div>
            </div>
    
            <table
              align="center"
              style="
                text-align: center;
                vertical-align: top;
                width: 600px;
                max-width: 600px;
                background-color: #ffffff;
              "
              width="600"
            >
              <tbody style="color: #343434">
                <tr>
                  <td
                    style="
                      width: 596px;
                      vertical-align: top;
                      padding-left: 30px;
                      padding-right: 30px;
                      padding-top: 30px;
                      padding-bottom: 40px;
                    "
                    width="596"
                  >
                    <h1
                      style="
                        font-size: 22px;
                        line-height: 34px;
                        font-family: 'Helvetica', Arial, sans-serif;
                        font-weight: 600;
                        text-decoration: none;
                        color: #000000;
                      "
                    >
                      Hello ` + name + `, !
    
                    <p
                      style="
                        line-height: 24px;
                        font-weight: 400;
                        text-decoration: none;
                      "
                    >
                      Welcome to <strong>LinuxDiary 5.0</strong>. We are excited to have you on board. Your account has been successfully created.
                    </p>
    
                    <p>
                      Thank you for choosing LinuxDiary 5.0.
                    </p>
    
                    <p>
                      Regards,<br />
                      LinuxDiary 5.0 Team
                    </p>
                  </td>
                </tr>
              </tbody>
            </table>
    
            <table
              align="center"
              style="
                text-align: center;
                vertical-align: top;
                width: 600px;
                max-width: 600px;
                background-color: #ffffff;
              "
              width="600"
            >
              <tbody>
                <tr>
                  <td
                    style="
                      width: 596px;
                      vertical-align: top;
                      padding-left: 0;
                      padding-right: 0;
                    "
                  >
                    <img
                      style="
                        text-align: center;
                        border-top-left-radius: 30px;
                        border-bottom-right-radius: 30px;
                        margin-bottom: 5px;
                      "
                      alt="Logo"
                      src="https://res.cloudinary.com/dduur8qoo/image/upload/v1689771850/wlug_white_logo_page-0001_u8efnh.jpg"
                      align="center"
                      width="200"
                      height="120"
                      class="CToWUd"
                      data-bit="iit"
                    />
                  </td>
                </tr>
    
                <tr style="margin-bottom: 30px" align="center">
                  <td align="center">
                    <a
                      href="https://www.instagram.com/wcewlug/"
                      target="_blank"
                      data-saferedirecturl="https://www.google.com/url?q=https://www.instagram.com/wcewlug/&amp;source=gmail&amp;ust=1680976985984000&amp;usg=AOvVaw16ObtJOZ1hpw9644RZ4oMM"
                      style="margin: 0 12px"
                      ><img
                        src="https://res.cloudinary.com/dduur8qoo/image/upload/v1689773467/Instagram_vn7dni_kzulby.png"
                        class="CToWUd"
                        data-bit="iit"
                        height="30"
                        width="30"
                    /></a>
                    <a
                      href="https://twitter.com/wcewlug"
                      target="_blank"
                      data-saferedirecturl="https://www.google.com/url?q=https://twitter.com/wcewlug&amp;source=gmail&amp;ust=1680976985984000&amp;usg=AOvVaw1ypHRKREADjq_cn0IRD2po"
                      ><img
                        src="https://res.cloudinary.com/dduur8qoo/image/upload/v1689772239/twitter-icon-square-logo-108D17D373-seeklogo.com_tjkqmo.png"
                        class="CToWUd"
                        data-bit="iit"
                        height="30"
                        width="30"
                        style="border-radius: 5px"
                    /></a>
                    <a
                      href="https://linkedin.com/company/wlug-club"
                      target="_blank"
                      data-saferedirecturl="https://www.google.com/url?q=https://linkedin.com/company/wlug-club&amp;source=gmail&amp;ust=1680976985984000&amp;usg=AOvVaw0TDo2Akq1O-un9s_gRi70t"
                      style="margin: 0 10px"
                      ><img
                        src="https://res.cloudinary.com/dduur8qoo/image/upload/v1685247353/linkedin_mg2ujv.png"
                        class="CToWUd"
                        data-bit="iit"
                        height="30"
                        width="30"
                        style="border-radius: 5px"
                    /></a>
                    <a
                      href="http://discord.wcewlug.org/join"
                      target="_blank"
                      data-saferedirecturl="https://www.google.com/url?q=http://discord.wcewlug.org/join&amp;source=gmail&amp;ust=1680976985984000&amp;usg=AOvVaw3PNiAyDSeiO1V36KVKeLZl"
                      style="margin: 0 1px"
                      ><img
                        src="https://res.cloudinary.com/dduur8qoo/image/upload/v1689771996/unnamed_m7lgs0.png"
                        class="CToWUd"
                        data-bit="iit"
                        height="30"
                        width="30"
                    /></a>
                    <a
                      href="https://www.youtube.com/@wlugwce"
                      target="_blank"
                      data-saferedirecturl="https://www.google.com/url?q=https://www.youtube.com/@wlugwce&amp;source=gmail&amp;ust=1680976985984000&amp;usg=AOvVaw1N-j63WkZ1tuHt3WgIqBsc"
                      style="margin: 0 10px"
                      ><img
                        src="https://res.cloudinary.com/dduur8qoo/image/upload/v1689773715/yt_logo_rgb_light_fwfqdx.png"
                        class="CToWUd"
                        data-bit="iit"
                        height="30"
                        width="30"
                    /></a>
                  </td>
                </tr>
              </tbody>
            </table>
    
            <table
              align="center"
              style="
                text-align: center;
                vertical-align: top;
                width: 600px;
                max-width: 600px;
                background-color: #ffffff;
              "
              width="600"
            >
              <tbody>
                <tr>
                  <td
                    style="
                      width: 596px;
                      vertical-align: top;
                      padding-top: 15px;
                      padding-left: 30px;
                      padding-right: 30px;
                    "
                    width="596"
                  >
                    <p
                      style="
                        font-size: 13px;
                        line-height: 18px;
                        font-family: 'Helvetica', Arial, sans-serif;
                        font-weight: 400;
                        text-decoration: none;
                        color: #000000;
                      "
                    >
                      Copyright &copy; 2024 WCE WLUG. All rights reserved. 
                      <br />
                      <br />
                      You are receiving this email because you opted in via our
                      website.
                      <br />
                      Our mailing address is:<br />
                      WCE, WLUG
                    </p>
                    <a
                      href="https://wcewlug.org"
                      target="_blank"
                      data-saferedirecturl="https://www.google.com/url?q=https://wcewlug.org&amp;source=gmail&amp;ust=1680976985984000&amp;usg=AOvVaw1h5YMuE2WETJts-iyzIb8t"
                      ><img
                        src="https://res.cloudinary.com/dduur8qoo/image/upload/v1689101481/favicon_qlhft1.png"
                        class="CToWUd"
                        data-bit="iit"
                        height="20"
                        width="20"
                        style="border-radius: 5px; float: left; margin-left: 50%"
                      />
                    </a>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </body>
    </html>`
}

func (u UserService) FileUpload(ctx context.Context, file multipart.File) (string, bool) {
	cld, _ := cloudinary.NewFromParams(os.Getenv("CLOUDINARY_CLOUD_NAME"), os.Getenv("CLOUDINARY_KEY"), os.Getenv("CLOUDINARY_SECRET"))

	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "LinuxDiary5.0",
	})

	if err != nil {
		return "", false
	}

	return uploadResult.SecureURL, true

}

func (u UserService) GetReferralLeaderboard(ctx context.Context, r *http.Request) ([]models.ReferralScore, error) {
	users, err := u.DbAdapter.GetUsers(ctx)

	leaderboard := make(map[string]int)

	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.ReferralCode != "" {
			leaderboard[user.ReferralCode]++
		}
	}

	var leaderboardArr []models.ReferralScore

	for key, value := range leaderboard {
		leaderboardArr = append(leaderboardArr, models.ReferralScore{ReferralCode: key, Score: value})
	}

	sort.Slice(leaderboardArr, func(i, j int) bool {
		return leaderboardArr[i].Score > leaderboardArr[j].Score
	})
	return leaderboardArr, nil
}
