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
    (2, "Bonsai", 5000000, 12, "https://modusaceh.co/files/images/20210923-gettyimages-598322873.jpg", "tanaman atau pohon yang dikerdilkan di dalam pot dangkal dengan tujuan membuat miniatur dari bentuk asli pohon besar yang sudah tua di alam bebas.");