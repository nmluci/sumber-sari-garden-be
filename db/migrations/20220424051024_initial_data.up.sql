INSERT INTO order_status(id, name) VALUES 
    (1, "Cart"),
    (2, "Not Paid"),
    (3, "Paid");

INSERT INTO user_role(id, name) VALUES
    (1, "admin"),
    (2, "customer");

INSERT INTO product_category(id, name) VALUES
    (1, "No Category"),
    (2, "Tree"),
    (3, "Flower"),
    (4, "Leaf");

INSERT INTO product(category_id, name, price, qty, url, description) VALUES
    (4, "Aglaonema", 30000, 199, "https://ik.imagekit.io/kompiangg/Sumber_Sari_Garden/Aglaonema_1_ySVV8li6w.jpg", "tanaman hias cocok untuk dekorasi taman"),
    (3, "Mawar", 16500, 123, "https://ik.imagekit.io/kompiangg/Sumber_Sari_Garden/Mawar_2_yfLXTpXwZT.jpg", "Bunga mawar Hidup sering ditanam di depan rumah dan menjadikan tampilan rumah memberikan kesan yang asri dan cantik."),
    (3, "Anggrek", 250000, 64, "https://ik.imagekit.io/kompiangg/Sumber_Sari_Garden/Anggrek_3_IaGLYqHeD.jpeg", "tanaman hias anggrek yang terkenal dengan keindahan karakteristiknya. Warna dari bunga anggrek bulan bermacam-macam, daun hijaunya lebar serta tangkai yang kehitaman. Tanaman hias ini mudah di temukan di Indonesia dan beberapa negara Asia lainnya, salah satunya adalah anggrek bulan Taiwan."),
    (3, "Adenium", 13000, 174, "https://images.tokopedia.net/img/cache/900/product-1/2020/7/6/batch-upload/batch-upload_de695634-9150-45df-8a4e-d300df1aac48.jpg", "Kemboja jepang atau adenium adalah spesies tanaman hias, batangnya besar, bagian bawahnya menyerupai umbi, batangnya tidak berkambium, akarnya dapat membesar menyerupai umbi, bentuk daunnya panjang ada yang lonjong, runcing, kecil, dan besar, warna bunganya bermacam-macam."),
    (4, "Lidah Mertua", 40000, 129, "https://images.tokopedia.net/img/cache/900/product-1/2019/9/11/6794510/6794510_612803b7-845d-47e8-901a-a3a8a083c12c_700_700.jpg", "Tanaman Sansivera Lidah mertua, bisa di tanam di indor/autdor, tinggi 30 cm s/d 50 cm"),
    (4, "Lili Paris", 37000, 38, "https://images.tokopedia.net/img/cache/900/VqbcmM/2020/11/13/8e68a3f9-ee98-4368-9c62-8e51b9403aa0.jpg", "Spider plants (Chlorophytum comosum) atau biasa dikenal juga dengan nama lili paris, dan untuk varian ini dikenal dengan nama lili paris ikal. Sesuai namanya, Lili paris ikal memiliki kumpulan daun seperti rumput yang melengkung."),
    (2, "Bambu Air", 50000, 48, "https://www.ruparupa.com/blog/wp-content/uploads/2022/01/bambu-hoki-pembawa-keberuntungan.jpeg", "Bambu air adalah spesies tanaman berbunga di keluarga Asparagaceae, asli Afrika Tengah. Itu dinamai tukang kebun Jerman-Inggris Henry Frederick Conrad Sander. Tanaman ini biasanya dipasarkan sebagai 'bambu keberuntungan'."),
    (4, "Peacy Lily", 13500, 342, "https://images.tokopedia.net/img/cache/900/product-1/2021/6/20/11672761/11672761_77d9dbe6-89cc-4f2f-9fe7-c2794527b302.jpg", "Peace lily atau dikenal sebagai lili perdamaian menjadi salah satu bunga penghias, baik di dalam maupun luar ruangan."),
    (4, "Rosemary", 12000, 948, "https://images.tokopedia.net/img/cache/900/VqbcmM/2021/8/12/96d0250a-26a3-4230-8ee0-3cf09aa71786.jpg", "Rosemary dikenal sebagai tanaman serbaguna karena bisa bisa digunakan sebagai rempah-rempah, minyak esensial, atau dibuat teh karena memiliki kandungan yang bermanfaat bagi kesehatan tubuh."),
    (4, "Palem", 32500, 329, "https://images.tokopedia.net/img/cache/900/VqbcmM/2020/10/19/d0c78831-b413-472f-b95c-7539ae826e6f.jpg", "Pohon palem ini adalah salah satu pembersih udara terbaik namanya adalah pohon palem kuning ini memiliki nama ilmiah Chrysalidocarpus lutescens atau Dypsis lutescens yang merupakan jenis pohon palem dengan cara merumpun."),
    (2, "Pohon", 130000, 391, "https://images.tokopedia.net/img/cache/900/VqbcmM/2022/4/14/30249fe3-0622-4515-a42b-d9f423b6aef4.png", "Kaktus merupakan nama yang diberikan untuk anggota tumbuhan berbunga famili Cactaceae yang biasa ditemukan di daerah-daerah kering. Memiliki ciri khas tampilan hijau berduri, kaktus dapat tumbuh dan bertahan hidup pada waktu yang lama tanpa air."),
    (2, "Bonsai", 5000000, 12, "https://modusaceh.co/files/images/20210923-gettyimages-598322873.jpg", "tanaman atau pohon yang dikerdilkan di dalam pot dangkal dengan tujuan membuat miniatur dari bentuk asli pohon besar yang sudah tua di alam bebas."),
    (3, "Kembang Sepatu", 23700, 30, "https://ik.imagekit.io/9tx59i8qwh/Photos/36b3c5f4-607d-4943-b7a1-467ffd4c9c19-2_eg8pG8she.jpg","Kembang sepatu adalah tanaman semak suku Malvaceae yang berasal dari Asia Timur dan banyak ditanam sebagai tanaman hias di daerah tropis dan subtropis. Bunga besar, berwarna merah dan tidak berbau."),
    (4, "Tradescantia Nanouk", 30000, 134, "https://ik.imagekit.io/9tx59i8qwh/Photos/38eb8ca0-e02a-43d5-a9c0-04fb7972bf6a_45zdY7Oq8.jpg", "Tradescantia Nanouk yang memiliki nama panjang Tradescantia albiflora 'Nanouk' atau di luar negeri sana dikenal juga dengan sebutan Fantasy Venice. Tanaman ini merupakan jenis tanaman yang memiliki warna daun sangat konsisten. Dalam 1 helai daun memiliki garis-garis merah muda, putih, ungu, dan hijau membentuk pola yang sangat indah"),
    (3, "Pink Mandevilla", 50000, 76, "https://ik.imagekit.io/9tx59i8qwh/Photos/bombshell_coral_pink_FVGBFpBFu.jpg", "tanaman rambat dengan bunga yang cantik cocok sebagai penghias pagar rumah"),
    (2, "Kaktus", 130000, 32, "https://ik.imagekit.io/9tx59i8qwh/Photos/kaktus_lKjl5-s68.png", "Kaktus merupakan nama yang diberikan untuk anggota tumbuhan berbunga famili Cactaceae yang biasa ditemukan di daerah-daerah kering. Memiliki ciri khas tampilan hijau berduri, kaktus dapat tumbuh dan bertahan hidup pada waktu yang lama tanpa air."),
    (3, "Teratai", 55000, 143, "https://ik.imagekit.io/9tx59i8qwh/Photos/teratai_zExd8UtlwB.jpg", "tanaman hias kolam ikan jenis tanaman teratai air. ready warna bunga putih, ungu, pink"),
    (3, "Petunia", 45000, 231, "https://ik.imagekit.io/9tx59i8qwh/Photos/petunia_IKi8RR1zL.jpg", "Petunia adalah suatu genus tumbuhan berbunga dari famili Solanaceae yang bunganya berbentuk trompet. Tumbuhan ini berasal dari Amerika Selatan. Secara fisik, tinggi tanaman ini antara 16-30 cm, bunganya ada yang bermahkota tunggal dan ada pula yang bermahkota ganda dengan warna yang bervariasi"),
    (4, "Begonia Polkadot", 30000, 510, "https://ik.imagekit.io/9tx59i8qwh/Photos/begonia_polkadot_tlp_wkXch.jpg", "Tanaman Begonia Polkadot Moka merupakan salah satu jenis tanaman perennial yang termasuk dalam famili Begoniaceae. Habitat asli tanaman begonia adalah daerah yang beriklim tropis dan subtropis. Tanaman ini bisa tumbuh subur di dataran rendah sampai dataran tinggi."),
    (4, "Keladi Putih", 30000, 102, "https://ik.imagekit.io/9tx59i8qwh/Photos/keladi_putih_0_VFqx4pL.png", "Daun lebar memanjang ke depan membentuk hati, daun tipis, corak warna yang mencolok seperti hijau, putih"),
    (4, "Kuping Gajah Mangkok", 80000, 53, "https://ik.imagekit.io/9tx59i8qwh/Photos/kuping_gajah_mangkok_AmSB-XkQp.jpg", "Berbeda dengan kuping gajah biasa, pada tanaman kuping gajah dorayaki warna dasar daun selain tulangnya adalah hijau keputihan, tidak hijau murni seperti kuping gajah hijau. Hijaunya juga tidak pekat."),
    (3, "Wijaya Kusuma", 25000, 32, "https://ik.imagekit.io/9tx59i8qwh/Photos/wijaya_kusuma_JP9OcpyBA.jpg", "Bunga Wijayakusuma atau disebut juga Bunga Wiku termasuk jenis tanaman kaktus yang mempunyai kelas dicotiledoneae. Tanaman ini berasal dari Mexico dan dapat hidup pada daerah dengan iklim sedang sampai beriklim tropis."),
    (2, "Asoka", 2500000, 11, "https://ik.imagekit.io/9tx59i8qwh/Photos/asokaa_OvjwZmjD-.jpg", "Pohon asoka adalah pohon yang dianggap suci oleh agama Hindu. Pohonnya akan mengeluarkan harum pada malam hari di bulan April dan Mei setiap tahunnya. Pohon tanaman ini sering diasosiasikan dengan cinta dan kesucian"),
    (4, "Monstera King Deliciosa", 350000, 21, "https://ik.imagekit.io/9tx59i8qwh/Photos/Monstera_King_Deliciosa_1PlAx1UFp.jpg", "Tanaman Monstera deliciosa king merupakan salah satu jenis monstera yang memiliki ukuran besar. Daun tanaman ini lebar dengan lubang-lubang di tengah hingga pinggiran daun. Semakin tua daunnya, maka warnanya akan semakin pekat."),
    (3, "Lavender", 12000, 75, "https://ik.imagekit.io/9tx59i8qwh/Photos/lavender_HsLY1us1Z.jpg", "Lavender atau lavendel atau Lavandula adalah genus tumbuhan berbunga dalam suku Lamiaceae yang tersusun atas 25-30 spesies. Asal tumbuhan ini adalah dari wilayah selatan Laut Tengah sampai Afrika tropis dan ke timur sampai India—Dunia Lama."),
    (4, "Glacier Ivy", 30000, 61, "https://ik.imagekit.io/9tx59i8qwh/Photos/glacier_ivy_yPjGRPPlZ.jpg", "Glacier ivy merupakan salah satu kultivar English ivy atau Common ivy yang memiliki daun bercorak varigata, tanaman merambat dari suku Araliaceae yang banyak dijumpai di Eropa dan Asia Barat. Daun glacier ivy berbentuk segitiga seperti jantung, merupakan daun varigata dengan bercak-bercak tak beraturan berwarna abu-abu dan hijau serta warna krem di bagian tepi daunnya.");