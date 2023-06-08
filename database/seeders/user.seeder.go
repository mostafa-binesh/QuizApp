package seeders

import (
	D "docker/database"
	M "docker/models"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func UserSeeder() {
	// cities := []string{
	// 	"تهران",
	// 	"اصفهان",
	// 	"شیراز",
	// 	"تبریز",
	// 	"مشهد",
	// 	"کرمان",
	// 	"اهواز",
	// 	"رشت",
	// 	"قم",
	// 	"کرمانشاه",
	// 	"بندرعباس",
	// 	"یزد",
	// 	"اصفهان",
	// 	"کاشان",
	// 	"ساری",
	// 	"خرم‌آباد",
	// 	"زاهدان",
	// 	"اراک",
	// 	"همدان",
	// 	"قزوین",
	// 	"کرج",
	// 	"سقز",
	// 	"سنندج",
	// 	"لاهیجان",
	// 	"مراغه",
	// 	"ملایر",
	// 	"بروجرد",
	// 	"خمینی شهر",
	// 	"گنبد کاووس",
	// 	"سراخس",
	// 	"قزوین",
	// 	"نیشابور",
	// 	"آمل",
	// 	"خوی",
	// 	"گرگان",
	// 	"برجند",
	// 	"مهاباد",
	// 	"قائم‌شهر",
	// 	"آبادان",
	// }
	names := []string{
		"سارا اسدی",
		"محمدرضا رحمانی",
		"نازنین عبادی",
		"پریسا علی‌نژاد",
		"ماهرخ احمدی",
		"سارینا داودی",
		"علیرضا جمشیدی",
		"آیدین کریمی",
		"نرجس موسوی",
		"محمدحسین صدرنژاد",
		"صدرا زمانی",
		"فاطمه خاکباز",
		"احمد معینی",
		"هانیه بحرینی",
		"رضوان امیری",
		"محمدرضا موحدی",
		"روژینا جعفری",
		"علیرضا علی‌پور",
		"پوریا پارسا",
		"مهدی خزاعی",
		"سعید محمدزاده",
		"شیدا مشایخی",
		"رضا حاتمی",
		"پروانه گلچین",
		"وحید ایرانمنش",
		"سمیرا پایدار",
		"علی رضا شهبازی",
		"ماهرخ شجاعی",
		"شیرین حسینی",
		"رضا محمدی",
		"نیما یزدانی",
		"زهرا عزیزی",
		"مهتاب مرادی",
		"ایمان مجتبوی",
		"محمدجواد محمودی",
		"آتنا کاوسی",
		"ایلیا احمدی",
		"رها فضائلی",
		"علی محمدی‌فر",
		"هادی رهبری",
		"آیلین محمدی",
		"لیلا رشیدی",
		"شروین شاهمرادی",
		"محمدرضا مهدوی",
		"نیایش خبازی",
		"نرگس جهانبخش",
		"احمدرضا صفرنژاد",
		"لیلا مهدوی",
		"شهیندخت دلاوری",
		"مهدی مستوفی",
	}
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator with current time
	for i := 0; i < 40; i++ {
		// cityRandomNumber := rand.Intn(len(cities))
		nameRandomNumber := rand.Intn(len(names))
		randomPhoneNumber := rand.Intn(100000000)                         // generate random number between 0 and 99999999
		randomIranPhoneNumber := fmt.Sprintf("09%09d", randomPhoneNumber) // print formatted string with leading zeros
		D.DB().Create(&M.User{
			Name: names[nameRandomNumber],
			// PhoneNumber:  fmt.Sprintf("%08d", rand.Intn(100000000)),
			PhoneNumber:  randomIranPhoneNumber,
			Password:     "This is my password",
			Role:         2,
			NationalCode: strconv.Itoa(rand.Intn(9000000000) + 1000000000), // Generate 10-digit number
			PersonalCode: strconv.Itoa(rand.Intn(9000000000) + 1000000000), // Generate 10-digit number
		})
	}
}
